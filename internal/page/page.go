package page

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// Page represents a single page.
type Page struct {
	sync.RWMutex
	Name         string `json:"name"`
	IsApp        bool   `json:"isApp"`
	sessions     map[string]*Session
	clients      map[*Client]bool
	clientsMutex sync.RWMutex
}

// NewPage creates a new instance of Page.
func NewPage(name string, isApp bool) *Page {
	p := &Page{}
	p.Name = name
	p.IsApp = isApp
	p.sessions = make(map[string]*Session)
	p.clients = make(map[*Client]bool)
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

func (p *Page) registerClient(client *Client) {
	p.clientsMutex.Lock()
	defer p.clientsMutex.Unlock()

	if _, ok := p.clients[client]; !ok {
		log.Printf("Registering %v client %s to %s",
			client.role, client.id, p.Name)

		p.clients[client] = true
	}
}

func (p *Page) unregisterClient(client *Client) {
	p.clientsMutex.Lock()
	defer p.clientsMutex.Unlock()

	log.Printf("Unregistering %v client %s from page %s",
		client.role, client.id, p.Name)

	delete(p.clients, client)
}
