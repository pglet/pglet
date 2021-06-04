package connection

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/pglet/pglet/internal/config"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type WebSocket struct {
	conn *websocket.Conn
	send chan []byte
	done chan bool
}

func NewWebSocket(conn *websocket.Conn) *WebSocket {
	cws := &WebSocket{
		conn: conn,
		send: make(chan []byte, 10),
		done: make(chan bool),
	}
	return cws
}

func (c *WebSocket) Start(handler ReadMessageHandler) {
	// start read/write loops
	go c.readLoop(handler)
	go c.writeLoop()
	<-c.done
}

func (c *WebSocket) Send(message []byte) {
	c.send <- message
}

func (c *WebSocket) readLoop(readHandler ReadMessageHandler) {
	defer func() {
		log.Println("Exiting WebSocket read loop")
		close(c.send)
		c.close()
	}()
	c.conn.SetReadLimit(int64(config.MaxWebSocketMessageSize()))
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Traceln("received pong")
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) || err == websocket.ErrReadLimit {
				log.Errorf("error: %v", err)
			}
			break
		}

		err = readHandler(message)
		if err != nil {
			log.Errorf("error processing WebSocket message: %v", err)
			break
		}
	}
}

func (c *WebSocket) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("Exiting WebSocket write loop")
		ticker.Stop()
		c.close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorf("Error creating WebSocket message writer: %v", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Errorf("Error writing WebSocket message: %v", err)
				return
			}

			if err := w.Close(); err != nil {
				log.Errorf("Error closing WebSocket message writer: %v", err)
				return
			}
		case <-ticker.C:
			log.Traceln("send ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("Error sending WebSocket PING message: %v", err)
				return
			}
		}
	}
}

func (c *WebSocket) close() {
	c.conn.Close()

	select {
	case c.done <- true:
	default:
		// no listeners
	}
}
