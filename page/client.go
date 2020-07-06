package page

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"

	// RegisterWebClientAction registers WS client as web (browser) client
	RegisterWebClientAction = "registerWebClient"

	// RegisterHostClientAction registers WS client as host (script) client
	RegisterHostClientAction = "registerHostClient"

	// PageCommandFromHostAction adds, sets, gets, disconnects or performs other page-related command from host
	PageCommandFromHostAction = "pageCommandFromHost"

	// PageEventFromWebAction receives click, change, expand/collapse and other events from browser
	PageEventFromWebAction = "pageEventFromWeb"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type ClientRole int

const (
	None ClientRole = iota
	WebClient
	HostClient
)

type Client struct {
	id      string
	role    ClientRole
	conn    *websocket.Conn
	session *Session
	send    chan []byte
}

type Message struct {
	ID      string          `json:"id"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type RegisterClientActionRequestPayload struct {
	PageName string `json:"pageName"`
	IsApp    bool   `json:"isApp"`
}

type RegisterClientActionResponsePayload struct {
	SessionID string `json:"sessionID"`
	Error     string `json:"error"`
}

type PageCommandActionRequestPayload struct {
	PageName  string `json:"pageName"`
	SessionID string `json:"sessionID"`
	Command   string `json:"command"`
}

type PageCommandActionResponsePayload struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

type PageEventActionPayload struct {
	PageName  string `json:"pageName"`
	SessionID string `json:"sessionID"`
	EventName string `json:"eventName"`
	EventData string `json:"eventData"`
}

type readPumpHandler = func(*Client, []byte) error

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoID() string {
	return uuid.New().String()
}

func (c *Client) readPump(readHandler readPumpHandler) {
	defer func() {
		c.session.unregisterClient(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		fmt.Println("received pong")
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
		err = readHandler(c, message)
		if err != nil {
			log.Printf("error processing message: %v", err)
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			fmt.Println("send ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		id:   autoID(),
		conn: conn,
		send: make(chan []byte, 256),
	}

	fmt.Printf("New Client %s is connected, total: %d\n", client.id, 0)

	// register client

	go client.readPump(readHandler)
	go client.writePump()
}

func readHandler(c *Client, message []byte) error {
	fmt.Printf("Message from %s: %v\n", c.id, string(message))

	// decode message
	msg := &Message{}
	err := json.Unmarshal(message, msg)
	if err != nil {
		return err
	}

	switch msg.Action {
	case RegisterWebClientAction:
		registerWebClient(c, msg)

	case RegisterHostClientAction:
		registerHostClient(c, msg)

	case PageCommandFromHostAction:
		executeCommandFromHostClient(c, msg)

	case PageEventFromWebAction:
		processPageEventFromWebClient(c, msg)
	}

	// echo back
	// time.Sleep(2 * time.Second)
	// c.send <- message

	return nil
}

func registerWebClient(client *Client, message *Message) {
	fmt.Println("Registering as web client")
	payload := new(RegisterClientActionRequestPayload)
	json.Unmarshal(message.Payload, payload)

	// assign client role
	client.role = WebClient

	// subscribe as host client
	page := Pages().Get(payload.PageName)

	var session *Session

	if !page.IsApp {
		// shared page
		// retrieve zero session
		session = page.sessions[ZeroSession]
	} else {
		// app page
		// TODO
	}

	// create page if not found
	session.registerClient(client)

	responsePayload, _ := json.Marshal(&RegisterClientActionResponsePayload{
		SessionID: session.ID,
		Error:     "",
	})

	response, _ := json.Marshal(&Message{
		ID:      message.ID,
		Payload: responsePayload,
	})

	client.send <- response
}

func registerHostClient(client *Client, message *Message) {
	fmt.Println("Registering as host client")
	payload := new(RegisterClientActionRequestPayload)
	json.Unmarshal(message.Payload, payload)

	// assign client role
	client.role = HostClient

	// subscribe as host client
	page := Pages().Get(payload.PageName)
	if page == nil {
		page = NewPage(payload.PageName, payload.IsApp)
		Pages().Add(page)
	}

	// retrieve zero session
	session := page.GetSession(ZeroSession)
	if session == nil {
		session = NewSession(page, ZeroSession)
		page.AddSession(session)
	}

	session.registerClient(client)

	responsePayload, _ := json.Marshal(&RegisterClientActionResponsePayload{
		SessionID: session.ID,
		Error:     "",
	})

	response, _ := json.Marshal(&Message{
		ID:      message.ID,
		Payload: responsePayload,
	})

	client.send <- response
}

func executeCommandFromHostClient(client *Client, message *Message) {
	fmt.Println("Page command from host client")

	payload := new(PageCommandActionRequestPayload)
	json.Unmarshal(message.Payload, payload)

	// process command
	// TODO
	fmt.Println("Command for page:", payload.PageName)

	// send response
	responsePayload, _ := json.Marshal(&PageCommandActionResponsePayload{
		Result: "Good",
		Error:  "",
	})

	response, _ := json.Marshal(&Message{
		ID:      message.ID,
		Payload: responsePayload,
	})

	client.send <- response

	// TODO
	// parse command
	// send response if parsing error
	// update page tree
	// if it's "shared" page - broadcast page change event to all web clients
	// if it's "app" page - send page change event to a web client with matching web client ID
}

func processPageEventFromWebClient(client *Client, message *Message) {
	fmt.Println("Page event from browser")
	// TODO
	// if it's "shared" page - broadcast event to all host clients
	// if it's "app" page - broadcast event to all host clients, event payload contains web client ID
}
