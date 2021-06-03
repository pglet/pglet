package connection

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

const (
	reconnectingAttempts = 30

	// Time allowed to write a message to the peer.
	//writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	//pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 60 * time.Second

	// Wait time for pong response
	pongTimeout = 5 * time.Second
)

type WebSocket struct {
	wsURL              string
	conn               *websocket.Conn
	send               chan []byte
	startReconnect     chan bool
	resumeReadLoop     chan bool
	resumeWriteLoop    chan bool
	terminateWriteLoop chan bool
	reconnectHandler   ReconnectHandler
	pongReceived       chan bool
}

func NewWebSocket(wsURL string) *WebSocket {
	cws := &WebSocket{
		wsURL:              wsURL,
		startReconnect:     make(chan bool),
		resumeReadLoop:     make(chan bool),
		resumeWriteLoop:    make(chan bool),
		terminateWriteLoop: make(chan bool),
		send:               make(chan []byte),
		pongReceived:       make(chan bool),
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

			select {
			case c.terminateWriteLoop <- true:
			default:
				// already terminated
			}

			c.reconnect()

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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Debugln("Exiting WebSocket write loop")
		ticker.Stop()
	}()
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
			c.reconnect()

			<-c.resumeWriteLoop
			log.Debugln("Resumed WebSocket write loop")

		case <-ticker.C:
			log.Debugln("Sending Ping")

			pongTimeout := time.NewTimer(pongTimeout)
			go func() {
				select {
				case <-pongTimeout.C:
					// re-connect
					log.Warnln("Pong receiving timeout")
					c.reconnect()
				case <-c.pongReceived:
					// cancel
					log.Println("Pong received")
				}
			}()

			//c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("Error sending WebSocket PING message: %v", err)
				return
			}
		case <-c.terminateWriteLoop:
			return
		}
	}
}

func (c *WebSocket) reconnect() {
	select {
	case c.startReconnect <- true:
	default:
		// reconnect is in progress
	}
}

func (c *WebSocket) reconnectLoop() {

	for {
		<-c.startReconnect

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

		if err == nil {
			c.conn.SetPongHandler(func(string) error {
				c.pongReceived <- true
				return nil
			})
			return
		}

		if attempt == totalAttempts {
			log.Printf("Failed to re-connect after %d attempts", totalAttempts)
			return
		}

		attempt += 1
		d := b.Duration()
		log.Printf("%s, reconnecting in %s", err, d)
		time.Sleep(d)
	}
}
