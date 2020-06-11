package page

import (
	"encoding/json"
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

func NewControlFromJson(jsonCtrl string) (Control, error) {
	ctrl := Control{}
	err := json.Unmarshal([]byte(jsonCtrl), &ctrl)
	if err != nil {
		return nil, err
	}
	return ctrl, nil
}

func (ctl Control) Id() string {
	return ctl["i"].(string)
}

func (ctl Control) ParentId() string {
	return ctl["p"].(string)
}

func (ctl Control) AddChildId(childId string) {
	childIds, _ := ctl["c"].([]string)
	// if !ok {
	// 	childIds = make([]string, 0, 1)
	// 	ctl["c"] = childIds
	// }
	ctl["c"] = append(childIds, childId)
}
