package connection

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

type WebSocket struct {
	sync.Mutex
	connected        bool
	wsURL            string
	conn             *websocket.Conn
	send             chan []byte
	stopWriteLoop    chan bool
	reconnectHandler ReconnectHandler
}

func NewWebSocket(wsURL string) *WebSocket {
	cws := &WebSocket{
		wsURL:         wsURL,
		send:          make(chan []byte),
		stopWriteLoop: make(chan bool),
	}
	return cws
}

func (c *WebSocket) Start(readHandler ReadMessageHandler, reconnectHandler ReconnectHandler) (err error) {

	c.reconnectHandler = reconnectHandler

	err = c.connect()
	if err != nil {
		return err
	}

	// start read/write loops
	go c.readLoop(readHandler)
	go c.writeLoop()
	return
}

func (c *WebSocket) Send(message []byte) {
	c.send <- message
}

func (c *WebSocket) readLoop(handler ReadMessageHandler) {
	log.Debugln("Starting WS read loop")
	for {
		_, bytesMessage, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorln("WS read error:", err)
			c.close()
			err = c.connect()
			if err != nil {
				log.Errorln("error re-connecting WS on read:", err)
				c.stopWriteLoop <- true
				return
			}
			log.Println("Re-connected WS on read")
			continue
		}

		handler(bytesMessage)
	}
}

func (c *WebSocket) writeLoop() {
	log.Debugln("Starting WS write loop")
	for {
		select {
		case message := <-c.send:
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err == nil {
				_, err = w.Write(message)
				if err == nil {
					err = w.Close()
					if err == nil {
						continue
					}
				}
			}

			log.Errorln("WS write error:", err)
			c.close()
			err = c.connect()
			if err != nil {
				log.Errorln("error re-connecting WS on write:", err)
				return
			}
			log.Println("Re-connected WS on write")
		case <-c.stopWriteLoop:
			return
		}
	}
}

func (c *WebSocket) connect() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.connected {
		return
	}

	b := &backoff.Backoff{
		Min:    1 * time.Second,
		Max:    5 * time.Minute,
		Jitter: true,
	}

	totalAttempts := 5

	for {
		log.Println("Connecting via WebSockets to:", c.wsURL)
		c.conn, _, err = websocket.DefaultDialer.Dial(c.wsURL, nil)

		if err == nil {
			c.connected = true
			return
		}

		totalAttempts -= 1
		if totalAttempts == 0 {
			return
		}

		d := b.Duration()
		log.Println("%s, reconnecting in %s", err, d)
		time.Sleep(d)
	}
}

func (c *WebSocket) close() {
	c.Lock()
	defer c.Unlock()
	c.connected = false

	if c.conn != nil {
		log.Println("Closing WS connection")
		c.conn.Close()
	}
}
