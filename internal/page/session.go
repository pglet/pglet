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
	// PageID is a reserved page ID
	PageID = "page"
)

type commandHandler = func(*Session, command.Command) (string, error)

var (
	commandHandlers = map[string]commandHandler{
		command.Add:     add,
		command.Addf:    add,
		command.Set:     set,
		command.Setf:    set,
		command.Get:     get,
		command.Clean:   clean,
		command.Cleanf:  clean,
		command.Remove:  remove,
		command.Removef: remove,
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
	s.AddControl(NewControl("page", "", PageID))
	s.clients = make(map[*Client]bool)
	return s
}

// ExecuteCommand executes command and returns the result
func (session *Session) ExecuteCommand(command command.Command) (result string, err error) {
	session.Lock()
	defer session.Unlock()

	log.Printf("Execute command for page %s session %s: %+v\n",
		session.Page.Name, session.ID, command)

	commandHandler := commandHandlers[command.Name]
	if commandHandler == nil {
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
	session.broadcastCommandToWebClients(msg)
	return id, nil
}

func get(session *Session, command command.Command) (result string, err error) {

	// command format must be:
	// get <control-id> <property>
	if len(command.Values) < 2 {
		return "", errors.New("'get' command should have control ID and property specified")
	}

	// control ID
	id := command.Values[0]

	ctrl, ok := session.Controls[id]
	if !ok {
		return "", fmt.Errorf("control with ID '%s' not found", id)
	}

	// control property
	prop := command.Values[1]

	v := ctrl.GetAttr(prop)

	if v == nil {
		return "", nil
	}

	return v.(string), nil
}

func set(session *Session, command command.Command) (result string, err error) {

	// command format must be:
	// get <control-id> <property>
	if len(command.Values) < 1 {
		return "", errors.New("'set' command should have control ID specified")
	}

	// control ID
	id := command.Values[0]

	ctrl, ok := session.Controls[id]
	if !ok {
		return "", fmt.Errorf("control with ID '%s' not found", id)
	}

	props := make(map[string]interface{})
	props["i"] = id

	// set control properties, except system ones
	for n, v := range command.Attrs {
		if !IsSystemAttr(n) {
			ctrl.SetAttr(n, v)
			props[n] = v
		}
	}

	payload := &UpdateControlPropsPayload{
		Props: make([]map[string]interface{}, 0, 0),
	}

	payload.Props = append(payload.Props, props)

	// broadcast control updates to all connected web clients
	session.broadcastCommandToWebClients(NewMessage(UpdateControlPropsAction, payload))
	return "", nil
}

func clean(session *Session, command command.Command) (result string, err error) {

	// command format must be:
	// clean <control-id>
	if len(command.Values) < 1 {
		return "", errors.New("'clean' command should have control ID specified")
	}

	// control ID
	id := command.Values[0]

	ctrl, ok := session.Controls[id]
	if !ok {
		return "", fmt.Errorf("control with ID '%s' not found", id)
	}

	session.cleanControl(ctrl)

	// output page
	// pJSON, _ := json.MarshalIndent(session.Controls, "", "  ")
	// log.Println(string(pJSON))

	// broadcast command to all connected web clients
	session.broadcastCommandToWebClients(NewMessage(CleanControlAction, &CleanControlPayload{
		ID: id,
	}))
	return "", nil
}

func remove(session *Session, command command.Command) (result string, err error) {

	// command format must be:
	// clean <control-id>
	if len(command.Values) < 1 {
		return "", errors.New("'clean' command should have control ID specified")
	}

	// control ID
	id := command.Values[0]

	ctrl, ok := session.Controls[id]
	if !ok {
		return "", fmt.Errorf("control with ID '%s' not found", id)
	}

	if ctrl.ParentID() == "" {
		return "", fmt.Errorf("root control '%s' cannot be deleted", id)
	}

	session.deleteControl(ctrl)

	// broadcast command to all connected web clients
	session.broadcastCommandToWebClients(NewMessage(RemoveControlAction, &RemoveControlPayload{
		ID: id,
	}))
	return "", nil
}

func (session *Session) UpdateControlProps(props []map[string]interface{}) {
	session.Lock()
	defer session.Unlock()

	for _, p := range props {
		id := p["i"].(string)
		if ctrl, ok := session.Controls[id]; ok {

			// patch control properties
			for n, v := range p {
				if !IsSystemAttr(n) {
					ctrl.SetAttr(n, v)
				}
			}
		}
	}
}

// NextControlID returns the next auto-generated control ID
func (session *Session) NextControlID() string {
	nextID := fmt.Sprintf("%s%d", ControlAutoIDPrefix, session.nextControlID)
	session.nextControlID++
	return nextID
}

// AddControl adds a control to a page
func (session *Session) AddControl(ctrl *Control) error {
	if _, exists := session.Controls[ctrl.ID()]; exists {
		return nil
	}
	session.Controls[ctrl.ID()] = ctrl

	// find parent
	parentID := ctrl.ParentID()
	if parentID != "" {
		parentctrl, ok := session.Controls[parentID]

		if !ok {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		parentctrl.AddChildID(ctrl.ID())
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

func (session *Session) cleanControl(ctrl *Control) {

	// delete all descendants
	for _, descID := range session.getAllDescendantIds(ctrl) {
		delete(session.Controls, descID)
	}

	// clean up children collection
	ctrl.RemoveChildren()
}

func (session *Session) deleteControl(ctrl *Control) {

	// delete all descendants
	for _, descID := range session.getAllDescendantIds(ctrl) {
		delete(session.Controls, descID)
	}

	// delete control itself
	delete(session.Controls, ctrl.ID())

	// remove control from parent's children collection
	session.Controls[ctrl.ParentID()].RemoveChild(ctrl.ID())
}

func (session *Session) getAllDescendantIds(ctrl *Control) []string {
	return session.getAllDescendantIdsRecursively(make([]string, 0, 0), ctrl.ID())
}

func (session *Session) getAllDescendantIdsRecursively(descendantIds []string, ID string) []string {
	ctrl := session.Controls[ID]
	childrenIds := ctrl.GetChildrenIds()
	result := append(descendantIds, childrenIds...)
	for _, childID := range childrenIds {
		result = append(result, session.getAllDescendantIdsRecursively(make([]string, 0, 0), childID)...)
	}
	return result
}

func getPageID() string {
	return PageID
}

func isAutoID(id string) bool {
	return id == PageID || strings.HasPrefix(id, ControlAutoIDPrefix)
}

func (session *Session) broadcastCommandToWebClients(msg *Message) {

	serializedMsg, _ := json.Marshal(msg)

	for c := range session.clients {
		if c.role == WebClient {
			c.send(serializedMsg)
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

	log.Printf("Unregistering %v client %s from session %s:%s",
		client.role, client.id, session.Page.Name, session.ID)

	delete(session.clients, client)
}
