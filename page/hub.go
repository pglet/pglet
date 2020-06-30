package page

import "sync"

// Hub connects the page with web and host clients
type Hub struct {
	sync.RWMutex
	page        *Page
	webClients  map[*Client]bool
	hostClients map[*Client]bool
}

// NewHub creates a new instance of page Hub.
func NewHub() (*Hub, error) {
	h := &Hub{}
	h.webClients = make(map[*Client]bool)
	h.hostClients = make(map[*Client]bool)
	return h, nil
}

func (p *Hub) RegisterHostClient(client *Client) {
	// TODO
}

func (p *Hub) RegisterWebClient(client *Client) {
	// TODO
}

func (p *Hub) UnregisterHostClient(client *Client) {
	// TODO
}

func (p *Hub) UnregisterWebClient(client *Client) {
	// TODO
}
