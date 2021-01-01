package store

import (
	"fmt"
	"log"
	"time"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/utils"
)

const (
	sessionIDKey              = "%d:%s"
	pageNextIDKey             = "page_next_id"               // Inc integer with the next page ID
	pageKey                   = "page:%s"                    // page JSON data
	pagesLastUpdatedKey       = "pages_last_updated"         // set of page names sorted by last updated Unix timestamp
	pageHostClientsKey        = "page_host_clients:%d"       // a Set with client IDs
	pageSessionsKey           = "page_sessions:%d"           // a Set with session IDs
	sessionKey                = "session:%d:%s"              // session JSON data
	sessionsLastUpdatedKey    = "sessions_last_updated"      // set of page:session IDs sorted by last updated Unix timestamp
	sessionNextControlIDField = "nextControlID"              // Inc integer with the next control ID for a given session
	sessionControlsKey        = "session_controls:%d:%s"     // session controls, value is JSON data
	sessionHostClientsKey     = "session_host_clients:%d:%s" // a Set with client IDs
	sessionWebClientsKey      = "session_web_clients:%d:%s"  // a Set with client IDs
)

//
// Pages
// ==============================

func GetPage(pageName string) *model.Page {
	var p model.Page
	cache.HashGetObject(fmt.Sprintf(pageKey, pageName), &p)
	if p.ID == 0 {
		return nil
	}
	return &p
}

func AddPage(page *model.Page) {

	// TODO - check if the page exists

	pageID := cache.Inc(pageNextIDKey, 1)
	page.ID = pageID
	cache.HashSet(fmt.Sprintf(pageKey, page.Name),
		"id", page.ID,
		"name", page.Name,
		"isApp", page.IsApp,
		"clientIP", page.ClientIP)
	SetPageLastUpdated(page)
}

func SetPageLastUpdated(page *model.Page) {
	cache.SortedSetAdd(pagesLastUpdatedKey, page.Name, time.Now().Unix())
}

func GetLastUpdatedPages(before int64) []string {
	return cache.SortedSetPopRange(pagesLastUpdatedKey, 0, before)
}

func DeletePage(pageName string) {
	page := GetPage(pageName)
	if page != nil {
		cache.Remove(fmt.Sprintf(pageKey, pageName))
		cache.SortedSetRemove(pagesLastUpdatedKey, pageName)
	}
}

//
// Page Host Clients
// ==============================

func GetPageSessions(page *model.Page) []string {
	return cache.SetGet(fmt.Sprintf(pageSessionsKey, page.ID))
}

func GetPageHostClients(page *model.Page) []string {
	return cache.SetGet(fmt.Sprintf(pageHostClientsKey, page.ID))
}

func AddPageHostClient(page *model.Page, clientID string) {
	cache.SetAdd(fmt.Sprintf(pageHostClientsKey, page.ID), clientID)
}

func RemovePageHostClient(page *model.Page, clientID string) {
	cache.SetRemove(fmt.Sprintf(pageHostClientsKey, page.ID), clientID)
}

//
// Sessions
// ==============================

func GetSession(page *model.Page, sessionID string) *model.Session {

	var session model.Session
	cache.HashGetObject(fmt.Sprintf(sessionKey, page.ID, sessionID), &session)
	if session.ID == "" {
		return nil
	}
	session.Page = page
	return &session
}

func AddSession(session *model.Session) {
	cache.HashSet(fmt.Sprintf(sessionKey, session.Page.ID, session.ID),
		"id", session.ID)
	cache.SetAdd(fmt.Sprintf(pageSessionsKey, session.Page.ID), session.ID)
	SetSessionLastUpdated(session)
}

func SetSessionLastUpdated(session *model.Session) {
	cache.SortedSetAdd(sessionsLastUpdatedKey, fmt.Sprintf(sessionIDKey, session.Page.ID, session.ID), time.Now().Unix())
}

func GetLastUpdatedSessions(before int64) []string {
	return cache.SortedSetPopRange(sessionsLastUpdatedKey, 0, before)
}

func DeleteSession(pageID int, sessionID string) {
	cache.SetRemove(fmt.Sprintf(pageSessionsKey, pageID), sessionID)
	cache.SortedSetRemove(sessionsLastUpdatedKey, fmt.Sprintf(sessionIDKey, pageID, sessionID))
	cache.Remove(fmt.Sprintf(sessionKey, pageID, sessionID))
	cache.Remove(fmt.Sprintf(sessionControlsKey, pageID, sessionID))
}

//
// Controls
// ==============================

func GetSessionNextControlID(session *model.Session) int {
	return cache.HashInc(fmt.Sprintf(sessionKey, session.Page.ID, session.ID), sessionNextControlIDField, 1)
}

func GetSessionControl(session *model.Session, ctrlID string) *model.Control {
	cj := cache.HashGet(fmt.Sprintf(sessionControlsKey, session.Page.ID, session.ID), ctrlID)
	if cj == "" {
		return nil
	}
	ctrl, err := model.NewControlFromJSON(cj)
	if err != nil {
		log.Fatal(err)
	}
	return ctrl
}

func GetAllSessionControls(session *model.Session) map[string]*model.Control {
	fields := cache.HashGetAll(fmt.Sprintf(sessionControlsKey, session.Page.ID, session.ID))
	controls := make(map[string]*model.Control, len(fields))
	for k, v := range fields {
		ctrl, _ := model.NewControlFromJSON(v)
		controls[k] = ctrl
	}
	return controls
}

func SetSessionControl(session *model.Session, ctrl *model.Control) {
	cj := utils.ToJSON(ctrl)
	cache.HashSet(fmt.Sprintf(sessionControlsKey, session.Page.ID, session.ID), ctrl.ID(), cj)
}

func DeleteSessionControl(session *model.Session, ctrlID string) {
	cache.HashRemove(fmt.Sprintf(sessionControlsKey, session.Page.ID, session.ID), ctrlID)
}

//
// Session Host Clients
// ==============================

func GetSessionHostClients(session *model.Session) []string {
	return cache.SetGet(fmt.Sprintf(sessionHostClientsKey, session.Page.ID, session.ID))
}

func AddSessionHostClient(session *model.Session, clientID string) {
	cache.SetAdd(fmt.Sprintf(sessionHostClientsKey, session.Page.ID, session.ID), clientID)
}

func RemoveSessionHostClient(session *model.Session, clientID string) {
	cache.SetRemove(fmt.Sprintf(sessionHostClientsKey, session.Page.ID, session.ID), clientID)
}

//
// Session Web Clients
// ==============================

func GetSessionWebClients(session *model.Session) []string {
	return cache.SetGet(fmt.Sprintf(sessionWebClientsKey, session.Page.ID, session.ID))
}

func AddSessionWebClient(session *model.Session, clientID string) {
	cache.SetAdd(fmt.Sprintf(sessionWebClientsKey, session.Page.ID, session.ID), clientID)
}

func RemoveSessionWebClient(session *model.Session, clientID string) {
	cache.SetRemove(fmt.Sprintf(sessionWebClientsKey, session.Page.ID, session.ID), clientID)
}
