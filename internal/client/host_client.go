package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/pglet/pglet/internal/client/connection"
	"github.com/pglet/pglet/internal/page"
)

type HostClient struct {

	// "ws" endpoint full URL
	wsURL string

	connectOnce sync.Once

	// active connection
	conn connection.Conn

	// pageSessionClients by "pageName:sessionID"
	pageSessionClients map[string]map[*PipeClient]bool
	pipeClientsMutex   sync.RWMutex

	// async calls registry
	calls map[string]chan *json.RawMessage

	// new page sessions
	newSessions map[string]chan string
	nsLock      sync.RWMutex
}

func NewHostClient(wsURL string) *HostClient {
	hc := &HostClient{}
	hc.wsURL = wsURL
	hc.pageSessionClients = make(map[string]map[*PipeClient]bool)
	hc.calls = make(map[string]chan *json.RawMessage)
	hc.newSessions = make(map[string]chan string)

	if wsURL == "" {
		// local/loopback connection
		hc.conn = connection.NewLocal()
	} else {
		// WebSocket connection
		hc.conn = connection.NewWebSocket(wsURL)
	}

	return hc
}

func (hc *HostClient) Start() (err error) {

	// connect only once
	hc.connectOnce.Do(func() {
		err = hc.conn.Start(hc.readHandler)
	})

	return
}

func (hc *HostClient) readHandler(bytesMessage []byte) (err error) {
	message := &page.Message{}
	err = json.Unmarshal(bytesMessage, message)
	if err == nil {

		//log.Println("Message to host client:", message)

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
	return
}

func (hc *HostClient) Call(ctx context.Context, action string, payload interface{}) *json.RawMessage {

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
	result := make(chan *json.RawMessage, 1)
	hc.calls[messageID] = result

	// send message
	hc.conn.Send(msg)

	// wait for result to arrive
	// TODO - implement timeout
	select {
	case <-ctx.Done():
		return nil
	case r := <-result:
		return r
	}
}

func (hc *HostClient) CallAndForget(action string, payload interface{}) {

	// serialize payload
	jsonPayload, _ := json.Marshal(payload)
	msg, _ := json.Marshal(&page.Message{
		Action:  action,
		Payload: jsonPayload,
	})

	// send message
	hc.conn.Send(msg)
}

func (hc *HostClient) RegisterPipeClient(pc *PipeClient) {
	hc.pipeClientsMutex.Lock()
	defer hc.pipeClientsMutex.Unlock()
	key := getPageSessionKey(pc.pageName, pc.sessionID)
	clients, ok := hc.pageSessionClients[key]
	if !ok {
		clients = make(map[*PipeClient]bool)
		hc.pageSessionClients[key] = clients
	}
	clients[pc] = true
}

func (hc *HostClient) UnregisterPipeClient(pc *PipeClient) {
	hc.pipeClientsMutex.Lock()
	defer hc.pipeClientsMutex.Unlock()
	key := getPageSessionKey(pc.pageName, pc.sessionID)
	clients, ok := hc.pageSessionClients[key]
	if ok {
		delete(clients, pc)
	}
}

func (hc *HostClient) broadcastPageEvent(rawPayload *json.RawMessage) error {
	// parse event
	payload := &page.PageEventPayload{}
	err := json.Unmarshal(*rawPayload, payload)
	if err != nil {
		return err
	}

	log.Println("Event:", payload)

	// iterate through all pipe clients
	key := getPageSessionKey(payload.PageName, payload.SessionID)
	clients, ok := hc.pageSessionClients[key]

	if ok {
		for client := range clients {
			eventMessage := fmt.Sprintf("%s %s %s",
				payload.EventTarget, payload.EventName, payload.EventData)
			client.emitEvent(eventMessage)
		}
	}

	return nil
}

func getPageSessionKey(pageName string, sessionID string) string {
	return fmt.Sprintf("%s:%s", pageName, sessionID)
}

func (hc *HostClient) notifySession(rawPayload *json.RawMessage) error {

	payload := new(page.SessionCreatedPayload)
	json.Unmarshal(*rawPayload, payload)

	log.Printf("Notify %s subscribers about new session %s\n", payload.PageName, payload.SessionID)
	select {
	case hc.PageNewSessions(payload.PageName) <- payload.SessionID:
		// Event sent to subscriber
	default:
		// No event listeners
	}

	return nil
}

func (hc *HostClient) PageNewSessions(pageName string) chan string {
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
