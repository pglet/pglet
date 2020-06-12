package page

import (
	"encoding/json"
)

// Control is an element of a page.
type Control map[string]interface{}

// NewControl initializes a new control object.
func NewControl(controlType string, parentID string, id string) Control {
	ctl := Control{}
	ctl["t"] = controlType
	ctl["i"] = id
	ctl["p"] = parentID
	return ctl
}

// NewControlFromJSON initializes a new control instance from JSON.
func NewControlFromJSON(jsonCtrl string) (Control, error) {
	ctrl := Control{}
	err := json.Unmarshal([]byte(jsonCtrl), &ctrl)
	if err != nil {
		return nil, err
	}
	return ctrl, nil
}

// ID returns control's ID.
func (ctl Control) ID() string {
	return ctl["i"].(string)
}

// ParentID returns the ID of parent control.
func (ctl Control) ParentID() string {
	return ctl["p"].(string)
}

// AddChildID appends the child to the parent control.
func (ctl Control) AddChildID(childID string) {
	childIds, _ := ctl["c"].([]string)
	// if !ok {
	// 	childIds = make([]string, 0, 1)
	// 	ctl["c"] = childIds
	// }
	ctl["c"] = append(childIds, childID)
}
