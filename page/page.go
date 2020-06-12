package page

import (
	"fmt"
	"sync"
)

// Page represents a single page.
type Page struct {
	sync.RWMutex
	Name     string             `json:"name"`
	Controls map[string]Control `json:"controls"`
}

// New creates a new instance of Page.
func New(name string) (*Page, error) {
	p := &Page{}
	p.Name = name
	p.Controls = make(map[string]Control)
	return p, Pages().Add(p)
}

// AddControl adds a control to a page
func (page *Page) AddControl(ctl Control) error {
	// find parent
	parentID := ctl.ParentID()
	if parentID != "" {
		page.RLock()
		parentCtl, ok := page.Controls[parentID]
		page.RUnlock()

		if !ok {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		page.Lock()
		parentCtl.AddChildID(ctl.ID())
		page.Unlock()
	}

	page.Lock()
	page.Controls[ctl.ID()] = ctl
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
