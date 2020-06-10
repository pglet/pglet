package pglet

import (
	"fmt"
)

type Control map[string]interface{}

// func (ctl Control) Id() (string, error) {
//   id, ok := ctl["id"].(string)
//   if !ok {
//     return "", errors.New("id must be a string")
//   }
//   return id, nil
// }

func NewControl(controlType string, parentId string, id string) Control {
	ctl := Control{}
	ctl["t"] = controlType
	ctl["i"] = id
	ctl["p"] = parentId
	return ctl
}

func (ctl Control) Id() string {
	return ctl["i"].(string)
}

func (ctl Control) ParentId() string {
	return ctl["p"].(string)
}

func (ctl Control) AddChildId(childId string) {
	childIds, ok := ctl["c"].([]string)
	if !ok {
		childIds = make([]string, 0, 1)
		ctl["c"] = childIds
	}
	ctl["c"] = append(childIds, childId)
}

func (page Page) AddControl(ctl Control) error {
	// find parent
	parentId := ctl.ParentId()
	if parentId != "" {
		parentCtl, ok := page.Controls[parentId]
		if !ok {
			return fmt.Errorf("Parent control with id '%s' not found.", parentId)
		}

		// update parent's childIds
		parentCtl.AddChildId(ctl.Id())
	}

	page.Controls[ctl.Id()] = ctl
	return nil
}
