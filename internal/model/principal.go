package model

type Principal struct {
	ID              int      `json:"id"`
	Username        string   `json:"username"`
	Email           string   `json:"email"`
	IsAuthenticated bool     `json:"isAuthenticated"`
	Groups          []string `json:"groups"`
}
