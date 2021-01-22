package page

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/page/command"
	"github.com/pglet/pglet/internal/pubsub"
	"github.com/pglet/pglet/internal/store"
	"github.com/pglet/pglet/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// ZeroSession is ID of zero session
	ZeroSession string = "0"
	// ControlAutoIDPrefix is a prefix for auto-generated control IDs
	ControlAutoIDPrefix = "_"
	// ControlIDSeparator is a symbol between parts of control ID
	ControlIDSeparator = ":"
	// ReservedPageID is a reserved page ID
	ReservedPageID = "page"
)

type commandHandlerFn = func(*command.Command) (string, error)

type sessionHandler struct {
	session  *model.Session
	handlers map[string]commandHandlerFn
}

func newSessionHandler(session *model.Session) sessionHandler {
	return sessionHandler{
		session: session,
	}
}

type addCommandBatchItem struct {
	command *command.Command
	control *model.Control
}

// NewSession creates a new instance of Session.
func newSession(page *model.Page, id string, clientIP string) *model.Session {
	s := &model.Session{}
	s.Page = page
	s.ID = id
	s.ClientIP = clientIP

	store.AddSession(s)

	h := newSessionHandler(s)
	h.addControl(model.NewControl("page", "", ReservedPageID))

	return s
}

func (h *sessionHandler) extendExpiration() {
	var expiresInMinutes int
	if h.session.ID == ZeroSession {
		expiresInMinutes = config.PageLifetimeMinutes()
	} else {
		expiresInMinutes = config.AppLifetimeMinutes()
	}
	store.SetSessionExpiration(h.session, time.Now().Add(time.Duration(expiresInMinutes)*time.Minute))
}

// ExecuteCommand executes command and returns the result
func (h *sessionHandler) execute(cmd *command.Command) (result string, err error) {
	sl := h.lockSession()
	defer sl.Unlock()

	log.Printf("Execute command for page %s session %s: %+v\n",
		h.session.Page.Name, h.session.ID, cmd)

	handlers := map[string]commandHandlerFn{
		command.Add:     h.add,
		command.Addf:    h.add,
		command.Set:     h.set,
		command.Setf:    h.set,
		command.Append:  h.appendHandler,
		command.Appendf: h.appendHandler,
		command.Get:     h.get,
		command.Clean:   h.clean,
		command.Cleanf:  h.clean,
		command.Remove:  h.remove,
		command.Removef: h.remove,
	}

	handler := handlers[strings.ToLower(cmd.Name)]
	if handler == nil {
		return "", fmt.Errorf("Unknown command: %s", cmd.Name)
	}

	return handler(cmd)
}

func (h *sessionHandler) add(cmd *command.Command) (result string, err error) {

	// parent ID
	topParentID := cmd.Attrs["to"]
	topParentAt := -1
	if ta, err := strconv.Atoi(cmd.Attrs["at"]); err == nil {
		topParentAt = ta
	}

	if topParentID == "" {
		topParentID = ReservedPageID
	}

	//log.Println("COMMAND:", utils.ToJSON(cmd))

	// "Add" commands to process
	batch := make([]*addCommandBatchItem, 0)

	// top command
	indent := 0
	if len(cmd.Values) > 0 {
		// single command
		batch = append(batch, &addCommandBatchItem{
			command: cmd,
		})
		indent = 2
	}

	// sub-commands
	for _, line := range cmd.Lines {
		if utils.WhiteSpaceOnly(line) {
			continue
		}

		childCmd, err := command.Parse(line, false)
		if err != nil {
			return "", err
		}
		childCmd.Name = "add"
		childCmd.Indent += indent
		batch = append(batch, &addCommandBatchItem{
			command: childCmd,
		})
	}

	// list of control IDs
	ids := make([]string, 0)

	// list of controls to broadcast
	payload := &AddPageControlsPayload{
		Controls: make([]*model.Control, 0),
	}

	affectedParents := make(map[string]bool)

	// process batch
	for i, batchItem := range batch {

		// first value must be control type
		if len(batchItem.command.Values) == 0 {
			return "", errors.New("Control type is not specified")
		}

		controlType := batchItem.command.Values[0]

		// other values go to boolean properties
		if len(batchItem.command.Values) > 1 {
			for _, v := range batchItem.command.Values[1:] {
				batchItem.command.Attrs[strings.ToLower(v)] = "true"
			}
		}

		parentID := ""
		parentAt := -1

		// find nearest parentID
		for pi := i - 1; pi >= 0; pi-- {
			if batch[pi].command.Indent < batchItem.command.Indent {
				parentID = batch[pi].control.ID()
				break
			}
		}

		// parent wasn't found - use the topmost one
		if parentID == "" {
			parentID = topParentID
			parentAt = topParentAt
		}

		// control ID
		id := batchItem.command.Attrs["id"]
		if id == "" {
			id = h.nextControlID()
		} else {
			// generate unique ID
			parentIDs := h.getControlParentIDs(parentID)
			id = strings.Join(append(parentIDs, id), ControlIDSeparator)
		}

		batchItem.control = model.NewControl(controlType, parentID, id)
		affectedParents[parentID] = true

		if parentAt != -1 {
			batchItem.control.SetAttr("at", parentAt)
			topParentAt++
		}

		for k, v := range batchItem.command.Attrs {
			if !model.IsSystemAttr(k) {
				batchItem.control.SetAttr(k, v)
			}
		}

		err = h.addControl(batchItem.control)
		if err != nil {
			return "", err
		}
		payload.Controls = append(payload.Controls, batchItem.control)
		ids = append(ids, id)
	}

	// re-read affected parents
	for i, ctrl := range payload.Controls {
		if affectedParents[ctrl.ID()] {
			payload.Controls[i] = store.GetSessionControl(h.session, ctrl.ID())
		}
	}

	//log.Println("CONTROLS:", utils.ToJSON(session.Controls))

	// broadcast new controls to all connected web clients
	h.broadcastCommandToWebClients(NewMessage(AddPageControlsAction, payload))
	return strings.Join(ids, " "), nil
}

func (h *sessionHandler) get(cmd *command.Command) (result string, err error) {

	// command format must be:
	// get <control-id> <property>
	if len(cmd.Values) < 2 {
		return "", errors.New("'get' command should have control ID and property specified")
	}

	// control ID
	id := cmd.Values[0]

	ctrl := h.getControl(id)
	if ctrl == nil {
		return "", fmt.Errorf("control with ID '%s' not found", id)
	}

	// control property
	prop := cmd.Values[1]

	v := ctrl.GetAttr(prop)

	if v == nil {
		return "", nil
	}

	return v.(string), nil
}

func (h *sessionHandler) set(cmd *command.Command) (result string, err error) {

	batch := make([]*command.Command, 0)

	// top command
	if len(cmd.Values) > 0 {
		// single command
		batch = append(batch, cmd)
	}

	// sub-commands
	for _, line := range cmd.Lines {
		if utils.WhiteSpaceOnly(line) {
			continue
		}

		childCmd, err := command.Parse(line, false)
		if err != nil {
			return "", err
		}
		childCmd.Name = "set"
		batch = append(batch, childCmd)
	}

	payload := &UpdateControlPropsPayload{
		Props: make([]map[string]interface{}, 0, 0),
	}

	for _, batchCmd := range batch {
		// command format must be:
		// get <control-id> <property>
		if len(batchCmd.Values) < 1 {
			return "", errors.New("'set' command should have control ID specified")
		}

		// control ID
		id := batchCmd.Values[0]

		ctrl := h.getControl(id)
		if ctrl == nil {
			return "", fmt.Errorf("control with ID '%s' not found", id)
		}

		// other values go to boolean properties
		if len(batchCmd.Values) > 1 {
			for _, v := range batchCmd.Values[1:] {
				batchCmd.Attrs[v] = "true"
			}
		}

		props := make(map[string]interface{})
		props["i"] = id

		// set control properties, except system ones
		for n, v := range batchCmd.Attrs {
			if !model.IsSystemAttr(n) {
				ctrl.SetAttr(n, v)
				props[n] = v
			}
		}
		err = store.SetSessionControl(h.session, ctrl)
		if err != nil {
			return "", err
		}

		payload.Props = append(payload.Props, props)
	}

	// broadcast control updates to all connected web clients
	h.broadcastCommandToWebClients(NewMessage(UpdateControlPropsAction, payload))
	return "", nil
}

func (h *sessionHandler) appendHandler(cmd *command.Command) (result string, err error) {

	batch := make([]*command.Command, 0)

	// top command
	if len(cmd.Values) > 0 {
		// single command
		batch = append(batch, cmd)
	}

	// sub-commands
	for _, line := range cmd.Lines {
		if utils.WhiteSpaceOnly(line) {
			continue
		}

		childCmd, err := command.Parse(line, false)
		if err != nil {
			return "", err
		}
		childCmd.Name = "append"
		batch = append(batch, childCmd)
	}

	payload := &AppendControlPropsPayload{
		Props: make([]map[string]string, 0, 0),
	}

	for _, batchCmd := range batch {
		// command format must be:
		// append control-id property=value property=value ...
		if len(batchCmd.Values) < 1 {
			return "", errors.New("'append' command should have control ID specified")
		}

		// control ID
		id := batchCmd.Values[0]

		ctrl := h.getControl(id)
		if ctrl == nil {
			return "", fmt.Errorf("control with ID '%s' not found", id)
		}

		props := make(map[string]string)
		props["i"] = id

		// set control properties, except system ones
		for n, v := range batchCmd.Attrs {
			if !model.IsSystemAttr(n) {
				ctrl.AppendAttr(n, v)
				props[n] = v
			}
		}
		err = store.SetSessionControl(h.session, ctrl)
		if err != nil {
			return "", err
		}

		payload.Props = append(payload.Props, props)
	}

	// broadcast control updates to all connected web clients
	h.broadcastCommandToWebClients(NewMessage(AppendControlPropsAction, payload))
	return "", nil
}

func (h *sessionHandler) clean(cmd *command.Command) (result string, err error) {

	// command format:
	//    clean [id_1] [id_2] ... [at=index]

	ids := make([]string, 0)
	if len(cmd.Values) == 0 {
		// clean page if no IDs specified
		ids = append(ids, ReservedPageID)
	} else {
		ids = append(ids, cmd.Values...)
	}

	at := -1
	if a, err := strconv.Atoi(cmd.Attrs["at"]); err == nil {
		at = a
	}

	if at != -1 && len(ids) > 1 {
		return "", errors.New("'at' cannot be specified with a list of IDs")
	}

	// control ID
	for i, id := range ids {
		ctrl := h.getControl(id)
		if ctrl == nil {
			return "", fmt.Errorf("control with ID '%s' not found", id)
		}

		if at != -1 {
			childIDs := ctrl.GetChildrenIds()
			if at > len(childIDs)-1 {
				return "", fmt.Errorf("'at' is out of range")
			}

			ids[i] = childIDs[at]
			ctrl = h.getControl(ids[i])
		}

		h.cleanControl(ctrl)
	}

	// broadcast command to all connected web clients
	h.broadcastCommandToWebClients(NewMessage(CleanControlAction, &CleanControlPayload{
		IDs: ids,
	}))
	return "", nil
}

func (h *sessionHandler) remove(cmd *command.Command) (result string, err error) {

	// command format:
	//    remove [id_1] [id_2] ... [at=index]

	at := -1
	if a, err := strconv.Atoi(cmd.Attrs["at"]); err == nil {
		at = a
	}

	ids := make([]string, 0)
	if len(cmd.Values) == 0 && at == -1 {
		return "", errors.New("'page' control cannot be removed")
	} else if len(cmd.Values) == 0 {
		ids = append(ids, ReservedPageID)
	} else {
		ids = append(ids, cmd.Values...)
	}

	if at != -1 && len(ids) > 1 {
		return "", errors.New("'at' cannot be specified with a list of IDs")
	}

	// control ID
	for i, id := range ids {
		ctrl := h.getControl(id)
		if ctrl == nil {
			return "", fmt.Errorf("control with ID '%s' not found", id)
		}

		if at != -1 {
			childIDs := ctrl.GetChildrenIds()
			if at > len(childIDs)-1 {
				return "", fmt.Errorf("'at' is out of range")
			}

			ids[i] = childIDs[at]
			ctrl = h.getControl(ids[i])
		}

		h.deleteControl(ctrl)
	}

	// broadcast command to all connected web clients
	h.broadcastCommandToWebClients(NewMessage(RemoveControlAction, &RemoveControlPayload{
		IDs: ids,
	}))
	return "", nil
}

func (h *sessionHandler) updateControlProps(props []map[string]interface{}) error {
	sl := h.lockSession()
	defer sl.Unlock()

	for _, p := range props {
		id := p["i"].(string)
		if ctrl := h.getControl(id); ctrl != nil {

			// patch control properties
			for n, v := range p {
				if !model.IsSystemAttr(n) {
					ctrl.SetAttr(n, v)
				}
			}

			err := store.SetSessionControl(h.session, ctrl)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// nextControlID returns the next auto-generated control ID
func (h *sessionHandler) nextControlID() string {
	return fmt.Sprintf("%s%d", ControlAutoIDPrefix, store.GetSessionNextControlID(h.session))
}

// addControl adds a control to a page
func (h *sessionHandler) addControl(ctrl *model.Control) error {
	if h.getControl(ctrl.ID()) != nil {
		return nil
	}
	err := store.SetSessionControl(h.session, ctrl)
	if err != nil {
		return err
	}

	// find parent
	parentID := ctrl.ParentID()
	if parentID != "" {
		parentctrl := h.getControl(parentID)

		if parentctrl == nil {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		if at := ctrl.At(); at != -1 {
			parentctrl.InsertChildID(ctrl.ID(), at)
		} else {
			parentctrl.AddChildID(ctrl.ID())
		}
		err = store.SetSessionControl(h.session, parentctrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *sessionHandler) getControlParentIDs(parentID string) []string {
	var result []string
	result = make([]string, 0)
	idParts := strings.Split(parentID, ControlIDSeparator)
	for _, idPart := range idParts {
		if !h.isAutoID(idPart) {
			result = append(result, idPart)
		}
	}
	return result
}

func (h *sessionHandler) cleanControl(ctrl *model.Control) {

	// delete all descendants
	for _, descID := range h.getAllDescendantIds(ctrl) {
		h.deleteSessionControl(descID)
	}

	// clean up children collection
	ctrl.RemoveChildren()
	store.SetSessionControl(h.session, ctrl)
}

func (h *sessionHandler) deleteControl(ctrl *model.Control) {

	// delete all descendants
	for _, descID := range h.getAllDescendantIds(ctrl) {
		h.deleteSessionControl(descID)
	}

	// delete control itself
	h.deleteSessionControl(ctrl.ID())

	// remove control from parent's children collection
	parentCtrl := h.getControl(ctrl.ParentID())
	parentCtrl.RemoveChild(ctrl.ID())
	store.SetSessionControl(h.session, parentCtrl)
}

func (h *sessionHandler) getAllDescendantIds(ctrl *model.Control) []string {
	return h.getAllDescendantIdsRecursively(make([]string, 0, 0), ctrl.ID())
}

func (h *sessionHandler) getAllDescendantIdsRecursively(descendantIds []string, ID string) []string {
	ctrl := h.getControl(ID)
	childrenIds := ctrl.GetChildrenIds()
	result := append(descendantIds, childrenIds...)
	for _, childID := range childrenIds {
		result = append(result, h.getAllDescendantIdsRecursively(make([]string, 0, 0), childID)...)
	}
	return result
}

func (h *sessionHandler) isAutoID(id string) bool {
	return id == ReservedPageID || strings.HasPrefix(id, ControlAutoIDPrefix)
}

func (h *sessionHandler) broadcastCommandToWebClients(msg *Message) {

	serializedMsg, _ := json.Marshal(msg)

	for _, clientID := range store.GetSessionWebClients(h.session) {
		pubsub.Send(clientChannelName(clientID), serializedMsg)
	}
}

func (h *sessionHandler) getControl(ctrlID string) *model.Control {
	return store.GetSessionControl(h.session, ctrlID)
}

func (h *sessionHandler) deleteSessionControl(ctrlID string) {
	store.DeleteSessionControl(h.session, ctrlID)
}

func (h *sessionHandler) lockSession() cache.Unlocker {
	return cache.Lock(fmt.Sprintf("session-lock-%s", h.session.ID))
}
