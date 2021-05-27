package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/pglet/pglet/internal/client/connection"
	"github.com/pglet/pglet/internal/page"
)

type PageRegistration struct {
	RegistrationRequest *page.RegisterHostClientRequestPayload
	Sessions            map[string]map[*PipeClient]bool
}

type HostClient struct {

	// "ws" endpoint full URL
	wsURL string

	connectOnce sync.Once

	// active connection
	conn connection.Conn

	// page registrations
	pages     map[string]*PageRegistration
	pagesLock sync.RWMutex

	// async calls registry
	calls map[string]chan *json.RawMessage

	// new page sessions
	newSessions map[string]chan string
	nsLock      sync.RWMutex
}

func NewHostClient(wsURL string) *HostClient {
	hc := &HostClient{}
	hc.wsURL = wsURL
	hc.pages = make(map[string]*PageRegistration)
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
		err = hc.conn.Start(hc.readHandler, hc.reconnectHandler)
	})

	return
}

func (hc *HostClient) readHandler(bytesMessage []byte) (err error) {

	//log.Debugln("Host client read message:", string(bytesMessage))

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
		log.Errorf("Unsupported message received: %s", bytesMessage)
	}
	return
}

func (hc *HostClient) reconnectHandler() error {
	return nil
}

func (hc *HostClient) RegisterPage(ctx context.Context, request *page.RegisterHostClientRequestPayload) (*page.RegisterHostClientResponsePayload, error) {
	// call server
	result := hc.Call(ctx, page.RegisterHostClientAction, request)

	// parse response
	response := &page.RegisterHostClientResponsePayload{}
	err := json.Unmarshal(*result, response)

	if err != nil {
		log.Errorln("Error parsing ConnectAppPage response:", err)
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	// update pages registry
	hc.pagesLock.Lock()
	defer hc.pagesLock.Unlock()

	pr, ok := hc.pages[response.PageName]
	if !ok {
		pr = &PageRegistration{
			Sessions: make(map[string]map[*PipeClient]bool),
		}
		hc.pages[response.PageName] = pr
	}
	pr.RegistrationRequest = request
	pr.RegistrationRequest.PageName = response.PageName

	return response, nil
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

func (hc *HostClient) RegisterPipeClient(pc *PipeClient) error {
	log.Debugf("Register pipe client for %s:%s\n", pc.pageName, pc.sessionID)

	hc.pagesLock.Lock()
	defer hc.pagesLock.Unlock()

	pr, ok := hc.pages[pc.pageName]
	if !ok {
		return fmt.Errorf("page or app %s is not registered", pc.pageName)
	}

	sessionClients, ok := pr.Sessions[pc.sessionID]
	if !ok {
		sessionClients = make(map[*PipeClient]bool)
		pr.Sessions[pc.sessionID] = sessionClients
	}
	sessionClients[pc] = true
	return nil
}

func (hc *HostClient) UnregisterPipeClient(pc *PipeClient) {
	log.Debugf("Unregister pipe client for %s:%s\n", pc.pageName, pc.sessionID)

	hc.pagesLock.Lock()
	defer hc.pagesLock.Unlock()

	pr, ok := hc.pages[pc.pageName]
	if !ok {
		return
	}

	sessionClients, ok := pr.Sessions[pc.sessionID]
	if ok {
		delete(sessionClients, pc)
	}
	if len(sessionClients) == 0 {
		delete(pr.Sessions, pc.sessionID)
	}

	if len(pr.Sessions) == 0 && !pr.RegistrationRequest.IsApp {
		delete(hc.pages, pc.pageName)
	}
}

func (hc *HostClient) broadcastPageEvent(rawPayload *json.RawMessage) error {
	// parse event
	payload := &page.PageEventPayload{}
	err := json.Unmarshal(*rawPayload, payload)
	if err != nil {
		return err
	}

	hc.pagesLock.RLock()
	defer hc.pagesLock.RUnlock()

	pr, ok := hc.pages[payload.PageName]
	if !ok {
		return nil
	}

	// iterate through all session pipe clients
	sessionClients, ok := pr.Sessions[payload.SessionID]
	if ok {
		for client := range sessionClients {
			eventMessage := fmt.Sprintf("%s %s %s",
				payload.EventTarget, payload.EventName, payload.EventData)
			client.emitEvent(eventMessage)
		}
	}

	return nil
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

func (hc *HostClient) CloseAppClients(pageName string) {
	log.Debugln("Closing inactive app clients", pageName)

	hc.pagesLock.Lock()
	defer hc.pagesLock.Unlock()

	pr, ok := hc.pages[pageName]
	if !ok {
		return
	}

	for _, clients := range pr.Sessions {
		for client := range clients {
			client.close()
		}
	}

	delete(hc.pages, pageName)
}

func (hc *HostClient) Close() {
	log.Debugf("Closing host client %s\n", hc.wsURL)

	hc.pagesLock.RLock()
	defer hc.pagesLock.RUnlock()

	for _, pr := range hc.pages {
		for _, clients := range pr.Sessions {
			for client := range clients {
				client.close()
			}
		}
	}
}
