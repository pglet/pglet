package page

import (
	"sync"
)

// Page represents a single page.
type Page struct {
	sync.RWMutex
	Name     string `json:"name"`
	IsApp    bool   `json:"isApp"`
	sessions map[string]*Session
}

// NewPage creates a new instance of Page.
func NewPage(name string, isApp bool) *Page {
	p := &Page{}
	p.Name = name
	p.IsApp = isApp
	p.sessions = make(map[string]*Session)
	return p
}

func (page *Page) GetSession(sessionID string) *Session {
	page.RLock()
	defer page.RUnlock()
	return page.sessions[sessionID]
}

func (page *Page) AddSession(s *Session) {
	page.Lock()
	defer page.Unlock()
	page.sessions[s.ID] = s
}
