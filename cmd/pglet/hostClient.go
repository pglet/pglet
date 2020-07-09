package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pglet/pglet/page"
)

type hostClient struct {

	// "ws" endpoint full URL
	wsURL string

	connectOnce sync.Once

	// active WebSocket connection
	conn *websocket.Conn

	// pageSessionClients by "pageName:sessionID"
	pageSessionClients map[string][]*pipeClient

	// async calls registry
	calls map[string]chan *json.RawMessage

	// new page sessions
	newSessions map[string]chan string
	nsLock      sync.RWMutex

	// send channel
	send chan []byte

	// used to break read/write loops
	done chan bool
}

func newHostClient(wsURL string) *hostClient {
	hc := &hostClient{}
	hc.wsURL = wsURL
	hc.pageSessionClients = make(map[string][]*pipeClient)
	hc.calls = make(map[string]chan *json.RawMessage)
	hc.newSessions = make(map[string]chan string)

	hc.send = make(chan []byte)
	hc.done = make(chan bool)
	return hc
}

func (hc *hostClient) start() (err error) {

	// connect only once
	hc.connectOnce.Do(func() {
		log.Println("Connecting to:", hc.wsURL)

		hc.conn, _, err = websocket.DefaultDialer.Dial(hc.wsURL, nil)

		if err != nil {
			return
		}

		go hc.readLoop()
		go hc.writeLoop()
	})

	return
}

func (hc *hostClient) readLoop() {
	defer close(hc.done)

	for {
		_, bytesMessage, err := hc.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		message := &page.Message{}
		err = json.Unmarshal(bytesMessage, message)
		if err == nil {
			if message.ID != "" {
				// this is callback message
				result, ok := hc.calls[message.ID]
				if ok {
					delete(hc.calls, message.ID)
					result <- &message.Payload
				}
			} else if message.Action == page.PageEventToHostAction {
				// event
				hc.broadcastPageEvent(&message.Payload)
			} else if message.Action == page.SessionCreatedAction {
				// new session
				hc.notifySession(&message.Payload)
			}
		} else {
			log.Printf("Unsupported message received: %s", bytesMessage)
		}
	}
}

func (hc *hostClient) writeLoop() {
	for {
		select {
		case message, ok := <-hc.send:
			if !ok {
				err := hc.conn.Close()
				if err != nil {
					log.Fatalln(err)
				}
				return
			}

			w, err := hc.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(hc.send)
			for i := 0; i < n; i++ {
				w.Write(<-hc.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (hc *hostClient) call(action string, payload interface{}) *json.RawMessage {

	// assign unique ID to the message
	messageID := uuid.New().String()

	// serialize payload
	jsonPayload, _ := json.Marshal(payload)
	msg, _ := json.Marshal(&page.Message{
		ID:      messageID,
		Action:  action,
		Payload: jsonPayload,
	})

	// create and register result channel
	result := make(chan *json.RawMessage)
	hc.calls[messageID] = result

	// send message
	hc.send <- msg

	// wait for result to arrive
	// TODO - implement timeout
	return <-result
}

func (hc *hostClient) callAndForget(action string, payload interface{}) {

	// serialize payload
	jsonPayload, _ := json.Marshal(payload)
	msg, _ := json.Marshal(&page.Message{
		Action:  action,
		Payload: jsonPayload,
	})

	// send message
	hc.send <- msg
}

func (hc *hostClient) registerPipeClient(pc *pipeClient) {
	key := getPageSessionKey(pc.pageName, pc.sessionID)
	clients, ok := hc.pageSessionClients[key]
	if !ok {
		clients = make([]*pipeClient, 0, 1)
	}
	hc.pageSessionClients[key] = append(clients, pc)
}

func (hc *hostClient) broadcastPageEvent(rawPayload *json.RawMessage) error {
	// parse event
	payload := &page.PageEventPayload{}
	err := json.Unmarshal(*rawPayload, payload)
	if err != nil {
		return err
	}

	// iterate through all pipe clients
	key := getPageSessionKey(payload.PageName, payload.SessionID)
	clients, ok := hc.pageSessionClients[key]

	if ok {
		for _, client := range clients {
			client.emitEvent(fmt.Sprintf("%s %s %s",
				payload.EventTarget, payload.EventName, payload.EventData))
		}
	}

	return nil
}

func getPageSessionKey(pageName string, sessionID string) string {
	return fmt.Sprintf("%s:%s", pageName, sessionID)
}

func (hc *hostClient) notifySession(rawPayload *json.RawMessage) error {

	payload := new(page.SessionCreatedPayload)
	json.Unmarshal(*rawPayload, payload)

	log.Printf("Notify %s subscribers about new session %s\n", payload.PageName, payload.SessionID)
	select {
	case hc.pageNewSessions(payload.PageName) <- payload.SessionID:
		// Event sent to subscriber
	default:
		// No event listeners
	}

	return nil
}

func (hc *hostClient) pageNewSessions(pageName string) chan string {
	hc.nsLock.Lock()
	defer hc.nsLock.Unlock()

	var ns chan string
	ns, ok := hc.newSessions[pageName]
	if !ok {
		ns = make(chan string)
		hc.newSessions[pageName] = ns
	}
	return ns
}
