package model

import (
	"strconv"
	"strings"
)

type Session struct {
	Page *Page  `json:"-"`
	ID   string `json:"id" redis:"id"`
}

func ParseSessionID(fullSessionID string) (pageID int, sessionID string) {
	parts := strings.Split(fullSessionID, ":")
	if len(parts) == 2 {
		pageID, _ = strconv.Atoi(parts[0])
		sessionID = parts[1]
	}
	return
}
