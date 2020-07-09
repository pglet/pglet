package page

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
)

const (
	ZeroSession string = ""
)

type commandHandler = func(*Session, Command) (string, error)

var (
	commandHandlers = map[string]commandHandler{
		Add:    add,
		Addr:   add,
		Set:    set,
		Get:    get,
		Insert: insert,
		Clean:  clean,
		Remove: remove,
	}
)

// Session represents an instance of a page.
type Session struct {
	sync.RWMutex
	Page          *Page              `json:"page"`
	ID            string             `json:"id"`
	Controls      map[string]Control `json:"controls"`
	nextControlID int
	clients       map[*Client]bool
	clientsMutex  sync.RWMutex
}

// NewSession creates a new instance of Page.
func NewSession(page *Page, id string) *Session {
	s := &Session{}
	s.Page = page
	s.ID = id
	s.Controls = make(map[string]Control)
	s.AddControl(NewControl("Page", "", s.NextControlID()))
	s.clients = make(map[*Client]bool)
	return s
}

func (session *Session) ExecuteCommand(command Command) (result string, err error) {

	log.Printf("Execute command for page %s session %s: %+v\n",
		session.Page.Name, session.ID, command)

	commandHandler := commandHandlers[command.Name]
	if commandHandler != nil {
		return commandHandler(session, command)
	}

	// result = fmt.Sprintf("a\nb\n%+v", command)
	// time.Sleep(2 * time.Second)

	return
}

func add(session *Session, command Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	go session.broadcastCommandToWebClients(command)
	return "", nil
}

func set(session *Session, command Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	go session.broadcastCommandToWebClients(command)
	return "", nil
}

func get(session *Session, command Command) (result string, err error) {

	// TODO - implement command

	return "", nil
}

func insert(session *Session, command Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	go session.broadcastCommandToWebClients(command)
	return "", nil
}

func clean(session *Session, command Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	go session.broadcastCommandToWebClients(command)
	return "", nil
}

func remove(session *Session, command Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	go session.broadcastCommandToWebClients(command)
	return "", nil
}

// NextControlID returns the next auto-generated control ID
func (session *Session) NextControlID() string {
	session.Lock()
	defer session.Unlock()
	nextID := strconv.Itoa(session.nextControlID)
	session.nextControlID++
	return nextID
}

// AddControl adds a control to a page
func (session *Session) AddControl(ctl Control) error {
	// find parent
	parentID := ctl.ParentID()
	if parentID != "" {
		session.RLock()
		parentCtl, ok := session.Controls[parentID]
		session.RUnlock()

		if !ok {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		session.Lock()
		parentCtl.AddChildID(ctl.ID())
		session.Unlock()
	}

	session.Lock()
	defer session.Unlock()
	session.Controls[ctl.ID()] = ctl
	return nil
}

func (session *Session) broadcastCommandToWebClients(command Command) {

	msgPayload := &PageCommandRequestPayload{
		PageName:  session.Page.Name,
		SessionID: session.ID,
		Command:   command,
	}

	msgPayloadRaw, _ := json.Marshal(msgPayload)

	msg, _ := json.Marshal(&Message{
		Action:  PageCommandFromHostAction,
		Payload: msgPayloadRaw,
	})

	for c := range session.clients {
		if c.role == WebClient {
			c.send <- msg
		}
	}
}

func (session *Session) registerClient(client *Client) {
	session.clientsMutex.Lock()
	defer session.clientsMutex.Unlock()

	if _, ok := session.clients[client]; !ok {
		log.Printf("Registering %v client %s to %s:%s",
			client.role, client.id, session.Page.Name, session.ID)

		session.clients[client] = true
	}
}

func (session *Session) unregisterClient(client *Client) {
	session.clientsMutex.Lock()
	defer session.clientsMutex.Unlock()

	log.Printf("Unregistering %v client %s from %s:%s",
		client.role, client.id, session.Page.Name, session.ID)

	delete(session.clients, client)
}
