package page

import (
	"fmt"
	"sync"
)

type Page struct {
	sync.RWMutex                    //`json:"-:`
	Name         string             `json:"name"`
	Controls     map[string]Control `json:"controls"`
}

func New(name string) (*Page, error) {
	p := &Page{}
	p.Name = name
	p.Controls = make(map[string]Control)
	return p, Pages().AddPage(p)
}

func (page Page) AddControl(ctl Control) error {
	// find parent
	parentID := ctl.ParentId()
	if parentID != "" {
		page.RLock()
		parentCtl, ok := page.Controls[parentID]
		page.RUnlock()

		if !ok {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		page.Lock()
		parentCtl.AddChildId(ctl.Id())
		page.Unlock()
	}

	page.Lock()
	page.Controls[ctl.Id()] = ctl
	page.Unlock()
	return nil
}

// func (p Page) MarshalJSON() ([]byte, error) {
// 	var tmp struct {
// 		Name     string             `json:"name"`
// 		Controls map[string]Control `json:"controls1"`
// 	}
// 	tmp.Name = p.Name
// 	tmp.Controls = p.Controls
// 	return json.MarshalIndent(&tmp, "", "  ")
// }
