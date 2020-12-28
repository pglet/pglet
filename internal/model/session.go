package model

type Session struct {
	Page *Page  `json:"-"`
	ID   string `json:"id"`
}
