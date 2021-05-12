package page

import (
	"encoding/json"
	"errors"
	"fmt"
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
)

type ClientRole string

const (
	None                    ClientRole = "None"
	WebClient                          = "Web"
	HostClient                         = "Host"
	pageNotFoundMessage                = "Page not found or application is not running."
	inactiveAppMessage                 = "Application is inactive. Please try refreshing this page later."
	loginRequiredMessage               = "login_required"
	clientExpirationSeconds            = 20
)

type Client struct {
	id           string
	role         ClientRole
	conn         connection.Conn
	clientIP     string
	subscription chan []byte
	sessions     map[string]*model.Session
	pages        map[string]*model.Page
	done         chan bool
}

func autoID() string {
	return uuid.New().String()
}

func NewClient(conn connection.Conn, clientIP string) *Client {
	c := &Client{
		id:       autoID(),
		conn:     conn,
		clientIP: clientIP,
		sessions: make(map[string]*model.Session),
		pages:    make(map[string]*model.Page),
		done:     make(chan bool, 1),
	}

	log.Println("Client IP:", clientIP)

	go c.subscribe()
	go c.extendExpiration()

	go func() {
		conn.Start(c.readHandler)
		c.done <- true
		c.unregister()
	}()

	log.Debugf("New Client %s is connected, total: %d\n", c.id, 0)

	return c
}

func (c *Client) subscribe() {
	c.subscription = pubsub.Subscribe(clientChannelName(c.id))
	for {
		select {
		case msg, more := <-c.subscription:
			if !more {
				return
			}
			c.send(msg)
		}
	}
}

func (c *Client) extendExpiration() {
	ticker := time.NewTicker(time.Duration(clientExpirationSeconds-2) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// extend client expiration
			store.SetClientExpiration(c.id, time.Now().Add(time.Duration(clientExpirationSeconds)*time.Second))

			// extend app session expiration
			if c.role == WebClient {
				for _, session := range c.sessions {
					h := newSessionHandler(session)
					h.extendExpiration()
				}
			}
		case <-c.done:
			return
		}
	}
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
	log.Println("Registering as web client")
	payload := new(RegisterWebClientRequestPayload)
	json.Unmarshal(message.Payload, payload)

	// assign client role
	c.role = WebClient

	// subscribe as web client
	page := store.GetPageByName(payload.PageName)

	response := &RegisterWebClientResponsePayload{
		Error: "",
	}

	if page == nil {
		response.Error = pageNotFoundMessage
	} else {
		var session *model.Session

		// check permissions
		if page.Permissions != "" {
			log.Debugln("Page permissions:", page.Permissions)
			response.Error = loginRequiredMessage
			response.LoginOptions = auth.GetLoginOptions(page.Permissions)
			goto response
		}

		if page.IsApp {
			// app page

			if len(store.GetPageHostClients(page)) == 0 {
				response.Error = inactiveAppMessage
				goto response
			}

			var sessionCreated bool
			if payload.SessionID != "" {
				// lookup for existing session
				session = store.GetSession(page, payload.SessionID)
			}

			// create new session
			if session == nil {
				if sessionsRateLimitReached(c.clientIP) {
					response.Error = fmt.Sprintf("A limit of %d new sessions per hour has been reached", config.LimitSessionsPerHour())
					goto response
				}

				session = newSession(page, uuid.New().String(), c.clientIP, payload.PageHash)
				sessionCreated = true
			} else {
				log.Debugf("Existing session %s found for %s page\n", session.ID, page.Name)
			}

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
				for _, clientID := range store.GetPageHostClients(page) {
					store.AddSessionHostClient(session.Page.ID, session.ID, clientID)
					pubsub.Send(clientChannelName(clientID), msg)
					break
				}

				log.Debugf("New session %s started for %s page\n", session.ID, page.Name)
			}

		} else {
			// shared page
			// retrieve zero session
			session = store.GetSession(page, ZeroSession)
			c.registerSession(session)

			log.Debugf("Connected to zero session of %s page\n", page.Name)
		}

		response.Session = &SessionPayload{
			ID:       session.ID,
			Controls: store.GetAllSessionControls(session),
		}
	}

response:

	responseMsg := NewMessageData(message.ID, RegisterWebClientAction, response)
	c.send(responseMsg)
}

func (c *Client) registerHostClient(message *Message) {
	responsePayload := &RegisterHostClientResponsePayload{
		SessionID: "",
		PageName:  "",
		Error:     "",
	}

	var err error
	var page *model.Page
	var pageName *model.PageName

	log.Println("Registering as host client")
	payload := new(RegisterHostClientRequestPayload)
	json.Unmarshal(message.Payload, payload)

	if !config.AllowRemoteHostClients() && c.clientIP != "" {
		err = fmt.Errorf("Remote host clients are not allowed")
		goto response
	} else if config.HostClientsAuthToken() != "" && config.HostClientsAuthToken() != payload.AuthToken {
		err = fmt.Errorf("Invalid auth token")
		goto response
	}

	// assign client role
	c.role = HostClient

	pageName, err = model.ParsePageName(payload.PageName)
	if err != nil {
		goto response
	}

	responsePayload.PageName = pageName.String()

	// retrieve page and then create if not exists
	page = store.GetPageByName(responsePayload.PageName)

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
		page = model.NewPage(responsePayload.PageName, payload.IsApp, payload.Permissions, c.clientIP)
		store.AddPage(page)
	}

	// make sure unath client has access to a given page
	if config.CheckPageIP() && page.ClientIP != c.clientIP {
		err = errors.New("Page name is already taken")
		goto response
	}

	if !page.IsApp {
		// retrieve zero session
		session := store.GetSession(page, ZeroSession)
		if session == nil {
			session = newSession(page, ZeroSession, c.clientIP, "")
		}
		c.registerSession(session)
		responsePayload.SessionID = session.ID
	} else {
		// register host client as an app server
		c.registerPage(page)
	}

response:

	if err != nil {
		responsePayload.Error = err.Error()
	}

	c.send(NewMessageData(message.ID, "", responsePayload))
}

func (c *Client) executeCommandFromHostClient(message *Message) {
	log.Debugln("Page command from host client")

	payload := new(PageCommandRequestPayload)
	json.Unmarshal(message.Payload, payload)

	responsePayload := &PageCommandResponsePayload{
		Result: "",
		Error:  "",
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

	// re-send events to all connected host clients
	//go func() {

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
	//}()

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
	}
}

func (c *Client) registerPage(page *model.Page) {

	log.Printf("Registering host client %s to handle '%s' sessions", c.id, page.Name)

	store.AddPageHostClient(page, c.id)
	c.pages[page.Name] = page
}

func (c *Client) unregisterPage(page *model.Page) {
	log.Printf("Unregistering host client %s from '%s' page", c.id, page.Name)

	store.RemovePageHostClient(page, c.id)
	delete(c.pages, page.Name)

	// delete app sessions
	if page.IsApp {
		go func() {
			clients := make([]string, 0)
			for _, sessionID := range store.GetPageHostClientSessions(page.ID, c.id) {
				session := &model.Session{
					Page: page,
					ID:   sessionID,
				}
				store.RemoveSessionHostClient(session.Page.ID, session.ID, c.id)

				sessionClients := store.GetSessionWebClients(session.Page.ID, session.ID)
				for _, clientID := range sessionClients {
					clients = append(clients, clientID)
					store.RemoveSessionWebClient(session.Page.ID, session.ID, clientID)
				}

				if _, ok := c.sessions[sessionID]; ok {
					delete(c.sessions, sessionID)
				}

				log.Debugln("Delete inactive app session:", page.ID, sessionID)
				store.DeleteSession(page.ID, sessionID)
			}

			store.RemovePageHostClientSessions(page.ID, c.id)

			if len(store.GetPageHostClients(page)) == 0 {
				store.DeletePage(page.ID)
			}

			for _, clientID := range clients {
				log.Debugln("Notify client which app become inactive:", clientID)

				msg := NewMessageData("", AppBecomeInactiveAction, &AppBecomeInactivePayload{
					Message: inactiveAppMessage,
				})
				pubsub.Send(clientChannelName(clientID), msg)
			}
		}()
	}

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

func (c *Client) unregister() {

	log.Printf("Unregistering client %s (%d sessions)", c.id, len(c.sessions))

	// unsubscribe from pubsub
	pubsub.Unsubscribe(c.subscription)

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
