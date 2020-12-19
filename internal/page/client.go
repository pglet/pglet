package page

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/pglet/pglet/internal/page/command"
	"github.com/pglet/pglet/internal/page/connection"
)

const (
	// RegisterWebClientAction registers WS client as web (browser) client
	RegisterWebClientAction = "registerWebClient"

	// RegisterHostClientAction registers WS client as host (script) client
	RegisterHostClientAction = "registerHostClient"

	// SessionCreatedAction notifies host clients about new sessions
	SessionCreatedAction = "sessionCreated"

	// PageCommandFromHostAction adds, sets, gets, disconnects or performs other page-related command from host
	PageCommandFromHostAction = "pageCommandFromHost"

	// PageEventFromWebAction receives click, change, expand/collapse and other events from browser
	PageEventFromWebAction = "pageEventFromWeb"

	// PageEventToHostAction redirects events from web to host clients
	PageEventToHostAction = "pageEventToHost"

	AddPageControlsAction = "addPageControls"

	UpdateControlPropsAction = "updateControlProps"

	AppendControlPropsAction = "appendControlProps"

	RemoveControlAction = "removeControl"

	CleanControlAction = "cleanControl"
)

type ClientRole string

const (
	None       ClientRole = "None"
	WebClient             = "Web"
	HostClient            = "Host"
)

type Client struct {
	id       string
	role     ClientRole
	conn     connection.Conn
	sessions map[*Session]bool
	pages    map[*Page]bool
}

type RegisterHostClientRequestPayload struct {
	PageName string `json:"pageName"`
	IsApp    bool   `json:"isApp"`
}

type RegisterHostClientResponsePayload struct {
	SessionID string `json:"sessionID"`
	PageName  string `json:"pageName"`
	Error     string `json:"error"`
}

type RegisterWebClientRequestPayload struct {
	PageName string `json:"pageName"`
	IsApp    bool   `json:"isApp"`
}

type RegisterWebClientResponsePayload struct {
	Session *Session `json:"session"`
	Error   string   `json:"error"`
}

type SessionCreatedPayload struct {
	PageName  string `json:"pageName"`
	SessionID string `json:"sessionID"`
}

type PageCommandRequestPayload struct {
	PageName  string          `json:"pageName"`
	SessionID string          `json:"sessionID"`
	Command   command.Command `json:"command"`
}

type PageCommandResponsePayload struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

type PageEventPayload struct {
	PageName    string `json:"pageName"`
	SessionID   string `json:"sessionID"`
	EventTarget string `json:"eventTarget"`
	EventName   string `json:"eventName"`
	EventData   string `json:"eventData"`
}

type AddPageControlsPayload struct {
	Controls []*Control `json:"controls"`
}

type UpdateControlPropsPayload struct {
	Props []map[string]interface{} `json:"props"`
}

type AppendControlPropsPayload struct {
	Props []map[string]string `json:"props"`
}

type RemoveControlPayload struct {
	IDs []string `json:"ids"`
}

type CleanControlPayload struct {
	IDs []string `json:"ids"`
}

func autoID() string {
	return uuid.New().String()
}

func NewClient(conn connection.Conn) *Client {
	c := &Client{
		id:       autoID(),
		conn:     conn,
		sessions: make(map[*Session]bool),
		pages:    make(map[*Page]bool),
	}

	go func() {
		conn.Start(c.readHandler)
		c.unregister()
	}()

	log.Printf("New Client %s is connected, total: %d\n", c.id, 0)

	return c
}

func (c *Client) readHandler(message []byte) error {
	log.Printf("Message from %s: %v\n", c.id, string(message))

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

	case PageEventFromWebAction:
		c.processPageEventFromWebClient(msg)

	case UpdateControlPropsAction:
		c.updateControlPropsFromWebClient(msg)
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
	page := Pages().Get(payload.PageName)

	response := &RegisterWebClientResponsePayload{
		Error: "",
	}

	if page == nil {
		response.Error = "Page not found or access denied"
	} else {
		var session *Session

		if !page.IsApp {
			// shared page
			// retrieve zero session
			session = page.sessions[ZeroSession]

			log.Printf("Connected to zero session of %s page\n", page.Name)
		} else {
			// app page
			// create new session
			session = NewSession(page, uuid.New().String())
			page.AddSession(session)

			log.Printf("New session %s started for %s page\n", session.ID, page.Name)
		}

		c.registerSession(session)

		if page.IsApp {
			// pick connected host client from page pool and notify about new session created
			sessionCreatedPayloadRaw, _ := json.Marshal(&SessionCreatedPayload{
				PageName:  page.Name,
				SessionID: session.ID,
			})

			msg, _ := json.Marshal(&Message{
				Action:  SessionCreatedAction,
				Payload: sessionCreatedPayloadRaw,
			})

			// TODO
			// pick first host client for now
			for c := range page.clients {
				if c.role == HostClient {
					c.registerSession(session)
					c.send(msg)
				}
			}
		}

		response.Session = session
	}

	responsePayload, _ := json.Marshal(response)

	responseMsg, _ := json.Marshal(&Message{
		ID:      message.ID,
		Action:  RegisterWebClientAction,
		Payload: responsePayload,
	})

	c.send(responseMsg)
}

func (c *Client) registerHostClient(message *Message) {
	log.Println("Registering as host client")
	payload := new(RegisterHostClientRequestPayload)
	json.Unmarshal(message.Payload, payload)

	responsePayload := &RegisterHostClientResponsePayload{
		SessionID: "",
		PageName:  "",
		Error:     "",
	}

	// assign client role
	c.role = HostClient

	pageName, err := parsePageName(payload.PageName)
	if err == nil {

		responsePayload.PageName = pageName.String()

		// retrieve page and then create if not exists
		page := Pages().Get(responsePayload.PageName)
		if page == nil {
			page = NewPage(responsePayload.PageName, payload.IsApp)
			Pages().Add(page)
		}

		if !page.IsApp {
			// retrieve zero session
			session := page.GetSession(ZeroSession)
			if session == nil {
				session = NewSession(page, ZeroSession)
				page.AddSession(session)
			}
			c.registerSession(session)
			responsePayload.SessionID = session.ID
		} else {
			// register host client as an app server
			c.registerPage(page)
		}
	} else {
		responsePayload.Error = err.Error()
	}

	responsePayloadRaw, _ := json.Marshal(responsePayload)

	response, _ := json.Marshal(&Message{
		ID:      message.ID,
		Payload: responsePayloadRaw,
	})

	c.send(response)
}

func (c *Client) executeCommandFromHostClient(message *Message) {
	log.Println("Page command from host client")

	payload := new(PageCommandRequestPayload)
	json.Unmarshal(message.Payload, payload)

	responsePayload := &PageCommandResponsePayload{
		Result: "",
		Error:  "",
	}

	// retrieve page and session
	page := Pages().Get(payload.PageName)
	if page != nil {
		session := page.GetSession(payload.SessionID)
		if session != nil {
			// process command
			result, err := session.ExecuteCommand(&payload.Command)
			responsePayload.Result = result
			if err != nil {
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
		responsePayloadRaw, _ := json.Marshal(responsePayload)

		response, _ := json.Marshal(&Message{
			ID:      message.ID,
			Payload: responsePayloadRaw,
		})

		c.send(response)
	}
}

func (client *Client) processPageEventFromWebClient(message *Message) {

	// web client can have only one session assigned
	var session *Session
	for s := range client.sessions {
		session = s
		break
	}

	log.Println("Page event from browser:", string(message.Payload),
		"PageName:", session.Page.Name, "SessionID:", session.ID)

	payload := new(PageEventPayload)
	json.Unmarshal(message.Payload, payload)

	// add page/session information to payload
	payload.PageName = session.Page.Name
	payload.SessionID = session.ID

	// message to host clients
	msgPayload, _ := json.Marshal(&payload)

	msg, _ := json.Marshal(&Message{
		Action:  PageEventToHostAction,
		Payload: msgPayload,
	})

	// re-send events to all connected host clients
	for c := range session.clients {
		if c.role == HostClient {
			c.send(msg)
		}
	}
}

func (client *Client) updateControlPropsFromWebClient(message *Message) {

	// web client can have only one session assigned
	var session *Session
	for s := range client.sessions {
		session = s
		break
	}

	payload := new(UpdateControlPropsPayload)
	json.Unmarshal(message.Payload, payload)

	log.Println("Update control props from web browser:", string(message.Payload),
		"PageName:", session.Page.Name, "SessionID:", session.ID, "Props:", payload.Props)

	log.Printf("%+v", payload.Props)

	// update control tree
	session.UpdateControlProps(payload.Props)

	// re-send the message to all connected web clients
	go func() {
		msg, _ := json.Marshal(message)

		for c := range session.clients {
			if c.role == WebClient && c.id != client.id {
				c.send(msg)
			}
		}
	}()
}

func (c *Client) registerPage(page *Page) {
	page.registerClient(c)
	c.pages[page] = true
}

func (c *Client) registerSession(session *Session) {
	session.registerClient(c)
	c.sessions[session] = true
}

func (client *Client) unregister() {

	// unregister from all sessions
	for session := range client.sessions {
		session.unregisterClient(client)
	}

	// unregister from all pages
	for page := range client.pages {
		page.unregisterClient(client)
	}
}
