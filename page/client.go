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

	// register WS client as web (browser) client
	ActionRegisterWebClientRequest = "registerWebClientRequest"

	// register WS client as host (script) client
	ActionRegisterHostClientRequest = "registerHostClientRequest"

	// add, set, get, disconnect or other page-related command from host
	ActionPageCommandFromHostRequest = "pageCommandFromHostRequest"

	// click, change, expand/collapse and other events from browser
	ActionPageEventFromWebRequest = "pageEventFromWebRequest"
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

type Client struct {
	id   string
	conn *websocket.Conn
	page *Page
	send chan []byte
}

type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type RegisterClientActionRequestPayload struct {
	PageName string `json:"pageName"`
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
		//c.hub.unregister <- c
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
	case ActionRegisterWebClientRequest:
		fmt.Println("Registering as web client")
		payload := new(RegisterClientActionRequestPayload)
		json.Unmarshal(msg.Payload, payload)

		// subscribe as web client
		page := Pages().Get(payload.PageName)
		page.RegisterWebClient(c)

	case ActionRegisterHostClientRequest:
		fmt.Println("Registering as host client")
		payload := new(RegisterClientActionRequestPayload)
		json.Unmarshal(msg.Payload, payload)

		// subscribe as host client
		page := Pages().Get(payload.PageName)
		page.RegisterHostClient(c)

	case ActionPageCommandFromHostRequest:
		fmt.Println("Page command from host client")
		// TODO

	case ActionPageEventFromWebRequest:
		fmt.Println("Page event from browser")
		// TODO
	}

	// echo back
	time.Sleep(2 * time.Second)
	c.send <- message

	return nil
}

func webClientHandler() {
	// read event (click, change, etc.)

	// send event to all subscribed page host clients
}

func hostClientHandler() {
	// read command

	// update page structure

	// write update to subscribed web clients
}
