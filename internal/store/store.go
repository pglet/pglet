package store

import (
	"fmt"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/model"
)

const (
	pageKey  = "page:%s"
	pagesKey = "pages"
)

func AddPage(page *model.Page) {

	// TODO - check if the page exists

	pageID := cache.Inc("page_next_id", 1)
	page.ID = pageID
	cache.SetObject(fmt.Sprintf(pageKey, page.Name), page, 0)
}

func GetPage(pageName string) *model.Page {
	p := new(model.Page)
	cache.GetObject(fmt.Sprintf(pageKey, pageName), &p)
	return p
}

func AddSession(page *model.Page, session *model.Session) {
	// TODO
}

func GetSession(page *model.Page, sessionID string) *model.Session {
	// TODO
	return nil
}
