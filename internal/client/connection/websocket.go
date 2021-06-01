package connection

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

const (
	reconnectingAttempts = 30
)

type WebSocket struct {
	wsURL              string
	conn               *websocket.Conn
	send               chan []byte
	reconnect          chan bool
	resumeReadLoop     chan bool
	resumeWriteLoop    chan bool
	terminateWriteLoop chan bool
	reconnectHandler   ReconnectHandler
}

func NewWebSocket(wsURL string) *WebSocket {
	cws := &WebSocket{
		wsURL:              wsURL,
		reconnect:          make(chan bool),
		resumeReadLoop:     make(chan bool),
		resumeWriteLoop:    make(chan bool),
		terminateWriteLoop: make(chan bool),
		send:               make(chan []byte),
	}
	return cws
}

func (c *WebSocket) Start(readHandler ReadMessageHandler, reconnectHandler ReconnectHandler) (err error) {

	c.reconnectHandler = reconnectHandler

	// initial connect
	err = c.connect(1)
	if err != nil {
		return
	}

	// start reconnect/read/write loops
	go c.reconnectLoop()
	go c.readLoop(readHandler)
	go c.writeLoop()
	return
}

func (c *WebSocket) Send(message []byte) {
	c.send <- message
}

func (c *WebSocket) readLoop(handler ReadMessageHandler) {
	log.Debugln("Starting WebSocket read loop")
	for {
		_, bytesMessage, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorln("WebSocket read error:", err)
			c.terminateWriteLoop <- true

			select {
			case c.reconnect <- true:
			default:
				// reconnect is in progress
			}

			<-c.resumeReadLoop
			log.Debugln("Resumed WebSocket read loop")
			go c.writeLoop()
			continue
		}

		handler(bytesMessage)
	}
}

func (c *WebSocket) writeLoop() {
	log.Debugln("Starting WebSocket write loop")
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

			log.Errorln("WebSocket write error:", err)

			select {
			case c.reconnect <- true:
			default:
				// reconnect is in progress
			}

			<-c.resumeWriteLoop
			log.Debugln("Resumed WebSocket write loop")

		case <-c.terminateWriteLoop:
			log.Debugln("Exiting WebSocket write loop")
			return
		}
	}
}

func (c *WebSocket) reconnectLoop() {

	for {
		<-c.reconnect

		if c.conn != nil {
			log.Println("Closing WebSocket connection")
			c.conn.Close()
		}

		err := c.connect(reconnectingAttempts)

		if err != nil {
			log.Errorf("Error reconnecting WebSocket: %s", err)
			return // TODO - what to do here?
		}

		log.Println("Re-connected WebSocket")

		select {
		case c.resumeReadLoop <- true:
		default:
			// no listeners
		}

		select {
		case c.resumeWriteLoop <- true:
		default:
			// no listeners
		}

		if c.reconnectHandler != nil {
			c.reconnectHandler(err == nil)
		}
	}
}

func (c *WebSocket) connect(totalAttempts int) (err error) {

	b := &backoff.Backoff{
		Min:    1 * time.Second,
		Max:    1 * time.Minute,
		Jitter: true,
	}

	attempt := 1

	for {
		log.Printf("Connecting via WebSockets to %s (attempt %d of %d)", c.wsURL, attempt, totalAttempts)
		c.conn, _, err = websocket.DefaultDialer.Dial(c.wsURL, nil)

		if err == nil || attempt == totalAttempts {
			return
		}

		attempt += 1
		d := b.Duration()
		log.Printf("%s, reconnecting in %s", err, d)
		time.Sleep(d)
	}
}
