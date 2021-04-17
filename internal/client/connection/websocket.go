package connection

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WebSocket struct {
	wsURL string
	conn  *websocket.Conn
	send  chan []byte
}

func NewWebSocket(wsURL string) *WebSocket {
	cws := &WebSocket{
		wsURL: wsURL,
		send:  make(chan []byte),
	}
	return cws
}

func (c *WebSocket) Start(handler ReadMessageHandler) (err error) {

	log.Println("Connecting via WebSockets to:", c.wsURL)
	c.conn, _, err = websocket.DefaultDialer.Dial(c.wsURL, nil)

	if err != nil {
		return
	}

	// start read/write loops
	go c.readLoop(handler)
	go c.writeLoop()
	return
}

func (c *WebSocket) Send(message []byte) {
	c.send <- message
}

func (hc *WebSocket) readLoop(handler ReadMessageHandler) {
	for {
		_, bytesMessage, err := hc.conn.ReadMessage()
		if err != nil {
			log.Errorln("read error:", err)
			return
		}

		handler(bytesMessage)
	}
}

func (c *WebSocket) writeLoop() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.Close()
				if err != nil {
					log.Fatalln(err)
				}
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
		}
	}
}
