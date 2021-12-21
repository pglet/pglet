package page

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/pglet/pglet/internal/auth"
	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/page/command"
	"github.com/pglet/pglet/internal/page/connection"
	"github.com/pglet/pglet/internal/pubsub"
	"github.com/pglet/pglet/internal/store"
	"github.com/pglet/pglet/internal/utils"
)

type ClientRole string

const (
	None                         ClientRole = "None"
	WebClient                               = "Web"
	HostClient                              = "Host"
	pageNotFoundMessage                     = "Page not found or application is not running."
	inactiveAppMessage                      = "Application is inactive. Please try refreshing this page later."
	signinRequiredMessage                   = "signin_required"
	clientRefreshIntervalSeconds            = 5
	clientExpirationSeconds                 = 20
)

type Client struct {
	id                   string
	role                 ClientRole
	conn                 connection.Conn
	clientIP             string
	principal            *auth.SecurityPrincipal
	subscription         chan []byte
	sessions             map[string]*model.Session
	pages                map[string]*model.Page
	exitSubscribe        chan bool
	exitExtendExpiration chan bool
}

func autoID() string {
	return uuid.New().String()
}

func NewClient(conn connection.Conn, clientIP string, principal *auth.SecurityPrincipal) *Client {
	c := &Client{
		id:                   autoID(),
		conn:                 conn,
		clientIP:             clientIP,
		principal:            principal,
		sessions:             make(map[string]*model.Session),
		pages:                make(map[string]*model.Page),
		exitExtendExpiration: make(chan bool),
	}

	go func() {
		normalClosure := conn.Start(c.readHandler)
		c.unregister(normalClosure)
	}()

	log.Printf("New WebSocket client connected from %s: %s", clientIP, c.id)

	return c
}

func (c *Client) register(role ClientRole) {

	if c.role != "" {
		return
	}

	log.Printf("Registering %s client: %s", role, c.id)

	c.role = role

	// subscribe PubSub
	c.subscription = pubsub.Subscribe(clientChannelName(c.id))
	go func() {
		for {
			select {
			case msg, more := <-c.subscription:
				if !more {
					log.Debugln("Exit subscribe():", c.id)
					return
				}
				c.send(msg)
			}
		}
	}()

	// run sliding expiration
	go func() {
		ticker := time.NewTicker(time.Duration(clientRefreshIntervalSeconds) * time.Second)
		defer ticker.Stop()
		for {
			// extend client expiration
			store.SetClientExpiration(c.id, time.Now().Add(time.Duration(clientExpirationSeconds)*time.Second))

			// extend app session expiration
			if c.role == WebClient {
				for _, session := range c.sessions {
					h := newSessionHandler(session)
					h.extendExpiration()
				}
			}

			select {
			case <-ticker.C:
			case <-c.exitExtendExpiration:
				log.Debugln("Exit extendExpiration():", c.id)
				return
			}
		}
	}()
}

func (c *Client) readHandler(message []byte) error {
	log.Debugf("Message from %s: %v\n", c.id, string(message))

	// decode message
	msg := &Message{}
	err := json.Unmarshal(message, msg)
	if err != nil {
		return err
	}

	switch msg.Action {
	case RegisterWebClientAction:
		c.registerWebClient(msg)

	case RegisterHostClientAction:
		c.registerHostClient(msg)

	case PageCommandFromHostAction:
		c.executeCommandFromHostClient(msg)

	case PageCommandsBatchFromHostAction:
		c.executeCommandsBatchFromHostClient(msg)

	case PageEventFromWebAction:
		c.processPageEventFromWebClient(msg)

	case UpdateControlPropsAction:
		c.updateControlPropsFromWebClient(msg)

	case InactiveAppFromHostAction:
		c.handleInactiveAppFromHostClient(msg)
	}

	return nil
}

func (c *Client) send(message []byte) {
	c.conn.Send(message)
}

func (c *Client) registerWebClient(message *Message) {
	log.Println("registerWebClient()")
	request := new(RegisterWebClientRequestPayload)
	json.Unmarshal(message.Payload, request)

	// subscribe as web client
	page := store.GetPageByName(request.PageName)

	response := &RegisterWebClientResponsePayload{
		Error: "",
	}

	if page == nil {
		response.Error = pageNotFoundMessage
	} else {
		var session *model.Session

		// check permissions
		if page.Permissions != "" && (c.principal == nil || !c.principal.HasPermissions(page.Permissions)) {
			log.Debugln("Required page permissions:", page.Permissions)
			response.Error = signinRequiredMessage
			response.SigninOptions = auth.GetSigninOptions(page.Permissions)
			goto response
		}

		if page.IsApp {
			// app page

			if len(store.GetPageHostClients(page.ID)) == 0 {
				response.Error = inactiveAppMessage
				goto response
			}

			var sessionCreated bool
			if request.SessionID != "" {
				sessionID, err := c.decryptSensitiveData(request.SessionID, c.clientIP)
				if err != nil {
					log.Errorf("error decrypting request.SessionID %s from %s: %s", request.SessionID, c.clientIP, err)
				}
				// lookup for existing session
				session = store.GetSession(page, sessionID)
			}

			// create new session
			if session == nil {
				if sessionsRateLimitReached(c.clientIP) {
					response.Error = fmt.Sprintf("A limit of %d new sessions per hour has been reached", config.LimitSessionsPerHour())
					goto response
				}

				session = newSession(page, uuid.New().String(), c.clientIP,
					request.PageHash, request.PageWidth, request.PageHeight)
				sessionCreated = true
			} else {
				log.Debugf("Existing session %s found for %s page\n", session.ID, page.Name)
			}

			c.register(WebClient)
			c.registerSession(session)

			if sessionCreated {

				// pick connected host client from page pool and notify about new session created
				msg := NewMessageData("", SessionCreatedAction, &SessionCreatedPayload{
					PageName:  page.Name,
					SessionID: session.ID,
				})

				// TODO
				// pick first host client for now
				// in the future we will implement load distribution algorithm
				for _, clientID := range store.GetPageHostClients(page.ID) {
					store.AddSessionHostClient(session.Page.ID, session.ID, clientID)
					pubsub.Send(clientChannelName(clientID), msg)
					break
				}

				log.Debugf("New session %s started for %s page\n", session.ID, page.Name)
			}

			userProps := map[string]string{
				"userid":       "",
				"userlogin":    "",
				"username":     "",
				"useremail":    "",
				"userclientip": "",
			}

			principalID := ""
			if c.principal != nil {
				principalID = c.principal.UID

				userProps = map[string]string{
					"userid":       c.principal.ID,
					"userlogin":    c.principal.Login,
					"username":     c.principal.Name,
					"useremail":    c.principal.Email,
					"userclientip": c.principal.ClientIP,
				}
			}

			// update page's user details
			pctl := store.GetSessionControl(session, "page")
			for k, v := range userProps {
				pctl.SetAttr(k, v)
			}
			store.SetSessionControl(session, pctl)

			if session.PrincipalID != principalID {
				// update session's principal
				store.SetSessionPrincipalID(session, principalID)

				authEventName := "signout"

				if session.PrincipalID != "" {

					authEventName = "signin"

					// hide signin dialog
					pctl := store.GetSessionControl(session, "page")
					pctl.SetAttr("signin", "")
					store.SetSessionControl(session, pctl)
				}

				changeEventProps := map[string]interface{}{
					"i":      "page",
					"signin": "",
				}

				// inject user props
				for k, v := range userProps {
					changeEventProps[k] = v
				}

				eventData := []map[string]interface{}{
					changeEventProps,
				}

				data, _ := json.Marshal(eventData)
				msg := NewMessageData("", PageEventToHostAction, &PageEventPayload{
					PageName:    session.Page.Name,
					SessionID:   session.ID,
					EventTarget: "page",
					EventName:   "change",
					EventData:   string(data),
				})

				for _, clientID := range store.GetSessionHostClients(page.ID, session.ID) {
					pubsub.Send(clientChannelName(clientID), msg)
				}

				// fire "page.signin/signout" event
				msg = NewMessageData("", PageEventToHostAction, &PageEventPayload{
					PageName:    page.Name,
					SessionID:   session.ID,
					EventTarget: "page",
					EventName:   authEventName,
				})

				for _, clientID := range store.GetSessionHostClients(page.ID, session.ID) {
					pubsub.Send(clientChannelName(clientID), msg)
				}
			}

		} else {
			// shared page
			// retrieve zero session
			session = store.GetSession(page, ZeroSession)
			c.register(WebClient)
			c.registerSession(session)

			log.Debugf("Connected to zero session of %s page\n", page.Name)
		}

		sessionID, err := c.encryptSensitiveData(session.ID, c.clientIP)
		if err != nil {
			log.Errorf("error encrypting session.ID: %s", err)
		}

		response.Session = &SessionPayload{
			ID:       sessionID,
			Controls: store.GetAllSessionControls(session),
		}
	}

response:

	responseMsg := NewMessageData(message.ID, RegisterWebClientAction, response)
	c.send(responseMsg)
}

func (c *Client) registerHostClient(message *Message) {

	var err error
	var page *model.Page
	var pageName *model.PageName

	log.Println("registerHostClient()")
	request := new(RegisterHostClientRequestPayload)
	json.Unmarshal(message.Payload, request)

	if request.HostClientID != "" {
		hostClientID, err := c.decryptSensitiveData(request.HostClientID, c.clientIP)
		if err != nil {
			log.Errorf("error decrypting request.HostClientID %s from %s: %s", request.HostClientID, c.clientIP, err)
		}
		if hostClientID != "" && c.id != hostClientID {
			log.Printf("Updating host client ID to %s", hostClientID)
			c.id = hostClientID
		}
	}

	response := &RegisterHostClientResponsePayload{
		SessionID: "",
		PageName:  "",
		Error:     "",
	}

	if !config.AllowRemoteHostClients() && c.clientIP != "" && c.clientIP != "::1" && c.clientIP != "127.0.0.1" {
		err = fmt.Errorf("Remote host clients are not allowed")
		goto response
	} else if config.HostClientsAuthToken() != "" && config.HostClientsAuthToken() != request.AuthToken {
		err = fmt.Errorf("Invalid auth token")
		goto response
	}

	pageName, err = model.ParsePageName(request.PageName)
	if err != nil {
		goto response
	}

	response.PageName = pageName.String()

	// retrieve page and then create if not exists
	page = store.GetPageByName(response.PageName)

	if page == nil {
		if pagesRateLimitReached(c.clientIP) {
			err = fmt.Errorf("A limit of %d new pages per hour has been reached", config.LimitPagesPerHour())
			goto response
		}

		// filter page name
		if pageName.IsReserved() {
			err = fmt.Errorf("Account or page name is reserved")
			goto response
		}

		// create new page
		page = model.NewPage(response.PageName, request.IsApp, request.Permissions, c.clientIP)
		store.AddPage(page)
	}

	// make sure unauth client has access to a given page
	if config.CheckPageIP() && page.ClientIP != c.clientIP {
		err = errors.New("Page name is already taken")
		goto response
	}

	// update page permissions
	if page.Permissions != request.Permissions {
		page.Permissions = request.Permissions
		store.UpdatePage(page)
	}

	// convert page to app
	if !page.IsApp && request.IsApp {
		page.IsApp = request.IsApp
		store.UpdatePage(page)

		// delete zero session
		store.DeleteSession(page.ID, ZeroSession)
	}

	c.register(HostClient)

	if !page.IsApp {
		// retrieve zero session
		session := store.GetSession(page, ZeroSession)
		if session == nil {
			session = newSession(page, ZeroSession, c.clientIP, "", "", "")
		}
		c.registerSession(session)
		response.SessionID = session.ID
	} else {
		// register host client as an app server
		c.registerPage(page)
	}

response:

	if err != nil {
		response.Error = err.Error()
	} else {
		hostClientID, err := c.encryptSensitiveData(c.id, c.clientIP)
		if err != nil {
			log.Errorf("error encrypting c.id: %s", err)
		}
		response.HostClientID = hostClientID
	}

	c.send(NewMessageData(message.ID, "", response))
}

func (c *Client) encryptSensitiveData(data string, clientIP string) (string, error) {
	result, err := utils.EncryptWithMasterKey([]byte(data + "|" + clientIP))
	if err != nil {
		return "", err
	}
	return utils.EncodeBase64(result), nil
}

func (c *Client) decryptSensitiveData(encrypted string, clientIP string) (string, error) {
	bytes, err := utils.DecodeBase64(encrypted)
	if err != nil {
		return "", err
	}
	plain, err := utils.DecryptWithMasterKey(bytes)
	if err != nil {
		return "", err
	}
	pair := strings.Split(string(plain), "|")
	if pair[1] != clientIP {
		return "", errors.New("IP address does not match")
	}
	return pair[0], nil
}

func (c *Client) executeCommandFromHostClient(message *Message) {
	log.Debugln("Page command from host client")

	payload := new(PageCommandRequestPayload)
	json.Unmarshal(message.Payload, payload)

	responsePayload := &PageCommandResponsePayload{
		Result: "",
		Error:  "",
	}

	if !payload.Command.IsSupported() {
		responsePayload.Error = fmt.Sprintf("unknown command: %s", payload.Command.Name)
		c.send(NewMessageData(message.ID, "", responsePayload))
		return
	}

	// retrieve page and session
	page := store.GetPageByName(payload.PageName)
	if page != nil {
		session := store.GetSession(page, payload.SessionID)
		if session != nil {
			// process command
			handler := newSessionHandler(session)
			result, err := handler.execute(payload.Command)
			responsePayload.Result = result

			if payload.Command.Name == command.Error {
				// session crashed on the client
				store.DeleteSession(page.ID, session.ID)
			} else if err != nil {
				handler.extendExpiration()
				responsePayload.Error = fmt.Sprint(err)
			}
		} else {
			responsePayload.Error = "Session not found or access denied"
		}
	} else {
		responsePayload.Error = "Page not found or access denied"
	}

	if payload.Command.ShouldReturn() {

		// send response
		c.send(NewMessageData(message.ID, "", responsePayload))
	}
}

func (c *Client) executeCommandsBatchFromHostClient(message *Message) {
	log.Debugln("Page commands batch from host client")

	payload := new(PageCommandsBatchRequestPayload)
	json.Unmarshal(message.Payload, payload)

	responsePayload := &PageCommandsBatchResponsePayload{
		Results: make([]string, 0),
		Error:   "",
	}

	for _, command := range payload.Commands {
		if !command.IsSupported() {
			responsePayload.Error = fmt.Sprintf("unknown command: %s", command.Name)
			c.send(NewMessageData(message.ID, "", responsePayload))
			return
		}
	}

	// retrieve page and session
	page := store.GetPageByName(payload.PageName)
	if page != nil {
		session := store.GetSession(page, payload.SessionID)
		if session != nil {
			// process command
			handler := newSessionHandler(session)
			results, err := handler.executeBatch(payload.Commands)
			responsePayload.Results = results
			if err != nil {
				handler.extendExpiration()
				responsePayload.Error = fmt.Sprint(err)
			}
		} else {
			responsePayload.Error = "Session not found or access denied"
		}
	} else {
		responsePayload.Error = "Page not found or access denied"
	}

	// send response
	c.send(NewMessageData(message.ID, "", responsePayload))
}

func (c *Client) processPageEventFromWebClient(message *Message) {

	// web client can have only one session assigned
	var session *model.Session
	for _, s := range c.sessions {
		session = s
		break
	}

	if session == nil {
		return
	}

	log.Debugln("Page event from browser:", string(message.Payload),
		"PageName:", session.Page.Name, "SessionID:", session.ID)

	payload := new(PageEventPayload)
	json.Unmarshal(message.Payload, payload)

	// add page/session information to payload
	payload.PageName = session.Page.Name
	payload.SessionID = session.ID

	msg := NewMessageData("", PageEventToHostAction, payload)

	// re-send events to all connected host clients
	for _, clientID := range store.GetSessionHostClients(session.Page.ID, session.ID) {
		pubsub.Send(clientChannelName(clientID), msg)
	}
}

func (c *Client) updateControlPropsFromWebClient(message *Message) error {

	// web client can have only one session assigned
	var session *model.Session
	for _, s := range c.sessions {
		session = s
		break
	}

	payload := new(UpdateControlPropsPayload)
	json.Unmarshal(message.Payload, payload)

	log.Debugln("Update control props from web browser:", string(message.Payload),
		"PageName:", session.Page.Name, "SessionID:", session.ID, "Props:", payload.Props)

	// update control tree
	handler := newSessionHandler(session)
	err := handler.updateControlProps(payload.Props)
	if err != nil {
		log.Errorln(err)
		return err
	}

	data, _ := json.Marshal(payload.Props)
	msg := NewMessageData("", PageEventToHostAction, &PageEventPayload{
		PageName:    session.Page.Name,
		SessionID:   session.ID,
		EventTarget: "page",
		EventName:   "change",
		EventData:   string(data),
	})

	for _, clientID := range store.GetSessionHostClients(session.Page.ID, session.ID) {
		pubsub.Send(clientChannelName(clientID), msg)
	}

	// re-send the message to all connected web clients
	go func() {
		msg, _ := json.Marshal(message)

		for _, clientID := range store.GetSessionWebClients(session.Page.ID, session.ID) {
			if clientID != c.id {
				pubsub.Send(clientChannelName(clientID), msg)
			}
		}
	}()
	return nil
}

func (c *Client) handleInactiveAppFromHostClient(message *Message) {
	payload := new(InactiveAppRequestPayload)
	json.Unmarshal(message.Payload, payload)

	log.Println("Handle inactive app from host client", payload.PageName)

	page, ok := c.pages[payload.PageName]
	if ok {
		c.unregisterPage(page)

		// delete app sessions
		go func() {
			webClients := make([]string, 0)
			for _, sessionID := range store.GetPageHostClientSessions(page.ID, c.id) {
				session := &model.Session{
					Page: page,
					ID:   sessionID,
				}
				store.RemoveSessionHostClient(session.Page.ID, session.ID, c.id)

				sessionClients := store.GetSessionWebClients(session.Page.ID, session.ID)
				for _, clientID := range sessionClients {
					webClients = append(webClients, clientID)
					store.RemoveSessionWebClient(session.Page.ID, session.ID, clientID)
				}

				if _, ok := c.sessions[sessionID]; ok {
					delete(c.sessions, sessionID)
				}

				log.Debugln("Delete inactive app session:", page.ID, sessionID)
				store.DeleteSession(page.ID, sessionID)
			}

			store.RemovePageHostClientSessions(page.ID, c.id)

			if len(store.GetPageHostClients(page.ID)) == 0 {
				store.DeletePage(page.ID)
			}

			notifyInactiveWebClients(webClients)
		}()
	}
}

func (c *Client) registerPage(page *model.Page) {

	log.Printf("Registering host client %s to handle '%s' sessions", c.id, page.Name)

	store.AddPageHostClient(page.ID, c.id)
	c.pages[page.Name] = page
}

func (c *Client) unregisterPage(page *model.Page) {
	log.Printf("Unregistering host client %s from '%s' page", c.id, page.Name)

	store.RemovePageHostClient(page.ID, c.id)
	delete(c.pages, page.Name)
}

func (c *Client) registerSession(session *model.Session) {

	log.Printf("Registering %v client %s to session %s:%s", c.role, c.id, session.Page.Name, session.ID)

	if c.role == WebClient {
		store.AddSessionWebClient(session.Page.ID, session.ID, c.id)
	} else {
		store.AddSessionHostClient(session.Page.ID, session.ID, c.id)
	}
	c.sessions[session.ID] = session

	h := newSessionHandler(session)
	h.extendExpiration()
}

func (c *Client) unregister(normalClosure bool) {

	log.Debugf("WebSocket client disconnected (normal closure=%t) from %s: %s", normalClosure, c.clientIP, c.id)

	if c.role == "" {
		return
	}

	log.Printf("Unregistering %s client %s (%d sessions)", c.role, c.id, len(c.sessions))

	// unsubscribe from pubsub
	pubsub.Unsubscribe(c.subscription)

	c.exitExtendExpiration <- true

	// unregister from all sessions
	for _, session := range c.sessions {
		log.Printf("Unregistering %v client %s from session %s:%s", c.role, c.id, session.Page.Name, session.ID)

		if c.role == WebClient {
			store.RemoveSessionWebClient(session.Page.ID, session.ID, c.id)
		} else {
			store.RemoveSessionHostClient(session.Page.ID, session.ID, c.id)
		}
	}

	// unregister from all pages
	for _, page := range c.pages {
		c.unregisterPage(page)
	}

	// expire client immediately
	if normalClosure {
		deleteExpiredClient(c.id)
	}
}

func pagesRateLimitReached(clientIP string) bool {
	limit := config.LimitPagesPerHour()
	if clientIP == "" || limit == 0 {
		return false
	}
	if cache.Inc(fmt.Sprintf("pages_limit:%s:%d", clientIP, time.Now().Hour()), 1, 1*time.Hour) > limit {
		return true
	}
	return false
}

func sessionsRateLimitReached(clientIP string) bool {
	limit := config.LimitSessionsPerHour()
	if clientIP == "" || limit == 0 {
		return false
	}
	if cache.Inc(fmt.Sprintf("sessions_limit:%s:%d", clientIP, time.Now().Hour()), 1, 1*time.Hour) > limit {
		return true
	}
	return false
}

func clientChannelName(clientID string) string {
	return fmt.Sprintf("client-%s", clientID)
}
