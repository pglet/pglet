package page

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/page/command"
)

const (
	// ZeroSession is ID of zero session
	ZeroSession string = ""
	// ControlAutoIDPrefix is a prefix for auto-generated control IDs
	ControlAutoIDPrefix = "_"
	// ControlIDSeparator is a symbol between parts of control ID
	ControlIDSeparator = ":"
)

type commandHandler = func(*Session, command.Command) (string, error)

var (
	commandHandlers = map[string]commandHandler{
		command.Add:    add,
		command.Addf:   add,
		command.Set:    set,
		command.Get:    get,
		command.Insert: insert,
		command.Clean:  clean,
		command.Remove: remove,
	}
)

// Session represents an instance of a page.
type Session struct {
	sync.RWMutex
	Page          *Page               `json:"page"`
	ID            string              `json:"id"`
	Controls      map[string]*Control `json:"controls"`
	nextControlID int
	clients       map[*Client]bool
	clientsMutex  sync.RWMutex
}

// NewSession creates a new instance of Page.
func NewSession(page *Page, id string) *Session {
	s := &Session{}
	s.Page = page
	s.ID = id
	s.Controls = make(map[string]*Control)
	s.AddControl(NewControl("page", "", s.NextControlID()))
	s.clients = make(map[*Client]bool)
	return s
}

// ExecuteCommand executes command and returns the result
func (session *Session) ExecuteCommand(command command.Command) (result string, err error) {

	log.Printf("Execute command for page %s session %s: %+v\n",
		session.Page.Name, session.ID, command)

	commandHandler := commandHandlers[command.Name]
	if commandHandler != nil {
		return "", fmt.Errorf("command '%s' does not have a handler", command.Name)
	}

	return commandHandler(session, command)
}

func add(session *Session, command command.Command) (result string, err error) {

	controlsFragment := command.Attrs["controls"]

	// first value must be control type
	if len(command.Values) == 0 && controlsFragment == "" {
		return "", errors.New("Control type is not specified")
	}

	controlType := command.Values[0]

	// parent ID
	parentID := command.Attrs["to"]

	if parentID == "" {
		parentID = getPageID()
	}

	// control ID
	id := command.Attrs["id"]
	if id == "" {
		id = session.NextControlID()
	} else {
		// generate unique ID
		parentIDs := getControlParentIDs(parentID)
		id = strings.Join(append(parentIDs, id), ControlIDSeparator)
	}

	ctrl := NewControl(controlType, parentID, id)

	for k, v := range command.Attrs {
		if !IsSystemAttr(k) {
			ctrl.SetAttr(k, v)
		}
	}

	session.AddControl(ctrl)

	// output page
	// pJSON, _ := json.MarshalIndent(session.Controls, "", "  ")
	// log.Println(string(pJSON))

	// update controls of all connected web cliens
	msg := NewMessage(AddPageControlsAction, &AddPageControlsPayload{
		Controls: []*Control{ctrl},
	})

	// broadcast command to all connected web clients
	go session.broadcastCommandToWebClients(msg)
	return id, nil
}

func set(session *Session, command command.Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	//go session.broadcastCommandToWebClients(command)
	return "", nil
}

func get(session *Session, command command.Command) (result string, err error) {

	// TODO - implement command

	return "", nil
}

func insert(session *Session, command command.Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	//go session.broadcastCommandToWebClients(command)
	return "", nil
}

func clean(session *Session, command command.Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	//go session.broadcastCommandToWebClients(command)
	return "", nil
}

func remove(session *Session, command command.Command) (result string, err error) {

	// TODO - implement command

	// broadcast command to all connected web clients
	//go session.broadcastCommandToWebClients(command)
	return "", nil
}

// NextControlID returns the next auto-generated control ID
func (session *Session) NextControlID() string {
	session.Lock()
	defer session.Unlock()
	nextID := fmt.Sprintf("%s%d", ControlAutoIDPrefix, session.nextControlID)
	session.nextControlID++
	return nextID
}

// AddControl adds a control to a page
func (session *Session) AddControl(ctl *Control) error {
	session.Lock()
	defer session.Unlock()
	if _, exists := session.Controls[ctl.ID()]; exists {
		return nil
	}
	session.Controls[ctl.ID()] = ctl

	// find parent
	parentID := ctl.ParentID()
	if parentID != "" {
		parentCtl, ok := session.Controls[parentID]

		if !ok {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		parentCtl.AddChildID(ctl.ID())
	}

	return nil
}

func getControlParentIDs(parentID string) []string {
	var result []string
	result = make([]string, 0)
	idParts := strings.Split(parentID, ControlIDSeparator)
	for _, idPart := range idParts {
		if !isAutoID(idPart) {
			result = append(result, idPart)
		}
	}
	return result
}

func getPageID() string {
	return fmt.Sprintf("%s%d", ControlAutoIDPrefix, 0)
}

func isAutoID(id string) bool {
	return strings.HasPrefix(id, ControlAutoIDPrefix)
}

func (session *Session) broadcastCommandToWebClients(msg *Message) {

	serializedMsg, _ := json.Marshal(msg)

	for c := range session.clients {
		if c.role == WebClient {
			c.send <- serializedMsg
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
