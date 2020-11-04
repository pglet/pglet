package page

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ConnWebSocket struct {
	conn *websocket.Conn
	send chan []byte
	done chan bool
}

func NewConnWebSocket(conn *websocket.Conn) *ConnWebSocket {
	cws := &ConnWebSocket{
		conn: conn,
		send: make(chan []byte),
		done: make(chan bool),
	}
	return cws
}

func (c *ConnWebSocket) Start(handler ReadMessageHandler) {
	// start read/write loops
	go c.readLoop(handler)
	go c.writeLoop()
	<-c.done
}

func (c *ConnWebSocket) Send(message []byte) {
	c.send <- message
}

func (c *ConnWebSocket) readLoop(readHandler ReadMessageHandler) {
	defer func() {
		c.done <- true
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Println("received pong")
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

		err = readHandler(message)
		if err != nil {
			log.Printf("error processing message: %v", err)
			break
		}
	}
}

func (c *ConnWebSocket) writeLoop() {
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
			log.Println("send ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
