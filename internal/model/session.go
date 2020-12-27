package model

type Session struct {
	Page *Page  `json:"page"`
	ID   string `json:"id"`
}

// NewSession creates a new instance of Page.
func NewSession(page *Page, id string) *Session {
	s := &Session{}
	s.Page = page
	s.ID = id
	return s
}
