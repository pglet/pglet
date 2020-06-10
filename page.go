package pglet

type Page struct {
	Name     string             `json:"name"`
	Controls map[string]Control `json:"controls"`
}
