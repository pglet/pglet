package model

// Page represents a single page.
type Page struct {
	ID       int    `json:"id" redis:"id"`
	Name     string `json:"name" redis:"name"`
	ClientIP string `json:"clientIP" redis:"clientIP"`
	IsApp    bool   `json:"isApp" redis:"isApp"`
	//sessions     map[string]*Session
	//clients      map[*Client]bool
	//clientsMutex sync.RWMutex
}

// NewPage creates a new instance of Page.
func NewPage(name string, isApp bool, clientIP string) *Page {
	p := &Page{}
	p.Name = name
	p.IsApp = isApp
	p.ClientIP = clientIP
	// p.sessions = make(map[string]*Session)
	// p.clients = make(map[*Client]bool)
	return p
}

// func (p *Page) registerClient(client *Client) {
// 	p.clientsMutex.Lock()
// 	defer p.clientsMutex.Unlock()

// 	if _, ok := p.clients[client]; !ok {
// 		log.Printf("Registering %v client %s to %s",
// 			client.role, client.id, p.Name)

// 		p.clients[client] = true
// 	}
// }

// func (p *Page) unregisterClient(client *Client) {
// 	p.clientsMutex.Lock()
// 	defer p.clientsMutex.Unlock()

// 	log.Printf("Unregistering %v client %s from page %s",
// 		client.role, client.id, p.Name)

// 	delete(p.clients, client)
// }
