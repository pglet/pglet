package store

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/cache"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/utils"
)

const (
	sessionIDKey              = "%d:%s"
	pageNextIDKey             = "page_next_id"                    // Inc integer with the next page ID
	pagesKey                  = "pages"                           // pages hash with pageName:pageID
	pageKey                   = "page:%d"                         // page data
	pageHostClientsKey        = "page_host_clients:%d"            // a Set with client IDs
	pageHostClientSessionsKey = "page_host_client_sessions:%d:%s" // a Set with sessionIDs
	pageSessionsKey           = "page_sessions:%d"                // a Set with session IDs
	clientSessionsKey         = "client_sessions:%s"              // a Set with session IDs
	sessionKey                = "session:%d:%s"                   // session data
	sessionsExpiredKey        = "sessions_expired"                // set of page:session IDs sorted by Unix timestamp of their expiration date
	clientsExpiredKey         = "clients_expired"                 // set of client IDs sorted by Unix timestamp of their expiration date
	sessionNextControlIDField = "nextControlID"                   // Inc integer with the next control ID for a given session
	sessionControlsKey        = "session_controls:%d:%s"          // session controls, value is JSON data
	sessionHostClientsKey     = "session_host_clients:%d:%s"      // a Set with client IDs
	sessionWebClientsKey      = "session_web_clients:%d:%s"       // a Set with client IDs
)

//
// Pages
// ==============================

func GetPageByName(pageName string) *model.Page {
	spid := cache.HashGet(pagesKey, pageName)
	if spid == "" {
		return nil
	}
	pageID, _ := strconv.Atoi(spid)
	return GetPageByID(pageID)
}

func GetPageByID(pageID int) *model.Page {
	var p model.Page
	cache.HashGetObject(fmt.Sprintf(pageKey, pageID), &p)
	if p.ID == 0 {
		return nil
	}
	return &p
}

func AddPage(page *model.Page) {

	// TODO - check if the page exists
	pageID := cache.Inc(pageNextIDKey, 1, 0)
	page.ID = pageID
	cache.HashSet(fmt.Sprintf(pageKey, page.ID),
		"id", page.ID,
		"name", page.Name,
		"isApp", page.IsApp,
		"clientIP", page.ClientIP)
	cache.HashSet(pagesKey, page.Name, page.ID)
}

func DeletePage(pageID int) {
	page := GetPageByID(pageID)
	if page == nil {
		log.Warnln("An attempt to delete inexisting page with ID", pageID)
		return
	}

	log.Println("Deleting page:", page.Name)
	for _, sessionID := range GetPageSessions(page.ID) {
		DeleteSession(page.ID, sessionID)
	}
	cache.Remove(fmt.Sprintf(pageHostClientsKey, page.ID))
	cache.Remove(fmt.Sprintf(pageKey, pageID))
	cache.Remove(pagesKey, page.Name)
}

//
// Page Host Clients
// ==============================

func GetPageSessions(pageID int) []string {
	return cache.SetGet(fmt.Sprintf(pageSessionsKey, pageID))
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
// Clients
// ==============================

func SetClientExpiration(clientID string, expires time.Time) {
	cache.SortedSetAdd(clientsExpiredKey, clientID, expires.Unix())
}

func GetExpiredClients() []string {
	return cache.SortedSetPopRange(clientsExpiredKey, 0, time.Now().Unix())
}

func DeleteExpiredClient(clientID string) {
	cache.SortedSetRemove(clientsExpiredKey, clientID)
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
}

func SetSessionExpiration(session *model.Session, expires time.Time) {
	cache.SortedSetAdd(sessionsExpiredKey, fmt.Sprintf(sessionIDKey, session.Page.ID, session.ID), expires.Unix())
}

func GetExpiredSessions() []string {
	return cache.SortedSetPopRange(sessionsExpiredKey, 0, time.Now().Unix())
}

func DeleteSession(pageID int, sessionID string) {
	cache.SetRemove(fmt.Sprintf(pageSessionsKey, pageID), sessionID)
	cache.SortedSetRemove(sessionsExpiredKey, fmt.Sprintf(sessionIDKey, pageID, sessionID))
	cache.Remove(fmt.Sprintf(sessionKey, pageID, sessionID))
	cache.Remove(fmt.Sprintf(sessionControlsKey, pageID, sessionID))
	cache.Remove(fmt.Sprintf(sessionHostClientsKey, pageID, sessionID))
	cache.Remove(fmt.Sprintf(sessionWebClientsKey, pageID, sessionID))
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

func SetSessionControl(session *model.Session, ctrl *model.Control) error {
	cj := utils.ToJSON(ctrl)
	success := cache.SetSessionControl(
		fmt.Sprintf(sessionKey, session.Page.ID, session.ID),
		fmt.Sprintf(sessionControlsKey, session.Page.ID, session.ID), ctrl.ID(), cj, config.LimitSessionSizeBytes())
	if !success {
		return fmt.Errorf("Session %d:%s size exceeds the maximum of %d bytes", session.Page.ID, session.ID, config.LimitSessionSizeBytes())
	}
	return nil
}

func DeleteSessionControl(session *model.Session, ctrlID string) {
	cache.RemoveSessionControl(
		fmt.Sprintf(sessionKey, session.Page.ID, session.ID),
		fmt.Sprintf(sessionControlsKey, session.Page.ID, session.ID), ctrlID)
}

//
// Session Host Clients
// ==============================

func GetSessionHostClients(session *model.Session) []string {
	return cache.SetGet(fmt.Sprintf(sessionHostClientsKey, session.Page.ID, session.ID))
}

func GetPageHostClientSessions(pageID int, clientID string) []string {
	return cache.SetGet(fmt.Sprintf(pageHostClientSessionsKey, pageID, clientID))
}

func AddSessionHostClient(session *model.Session, clientID string) {
	cache.SetAdd(fmt.Sprintf(sessionHostClientsKey, session.Page.ID, session.ID), clientID)
	cache.SetAdd(fmt.Sprintf(pageHostClientSessionsKey, session.Page.ID, clientID), session.ID)
	cache.SetAdd(fmt.Sprintf(clientSessionsKey, clientID), fmt.Sprintf(sessionIDKey, session.Page.ID, session.ID))
}

func RemoveSessionHostClient(session *model.Session, clientID string) {
	cache.SetRemove(fmt.Sprintf(sessionHostClientsKey, session.Page.ID, session.ID), clientID)
	cache.SetRemove(fmt.Sprintf(pageHostClientSessionsKey, session.Page.ID, clientID), session.ID)
	cache.SetRemove(fmt.Sprintf(clientSessionsKey, clientID), fmt.Sprintf(sessionIDKey, session.Page.ID, session.ID))
}

func RemovePageHostClientSessions(pageID int, clientID string) {
	cache.Remove(fmt.Sprintf(pageHostClientSessionsKey, pageID, clientID))
}

//
// Session Web Clients
// ==============================

func GetSessionWebClients(session *model.Session) []string {
	return cache.SetGet(fmt.Sprintf(sessionWebClientsKey, session.Page.ID, session.ID))
}

func AddSessionWebClient(session *model.Session, clientID string) {
	cache.SetAdd(fmt.Sprintf(sessionWebClientsKey, session.Page.ID, session.ID), clientID)
	cache.SetAdd(fmt.Sprintf(clientSessionsKey, clientID), fmt.Sprintf(sessionIDKey, session.Page.ID, session.ID))
}

func RemoveSessionWebClient(session *model.Session, clientID string) {
	cache.SetRemove(fmt.Sprintf(sessionWebClientsKey, session.Page.ID, session.ID), clientID)
	cache.SetRemove(fmt.Sprintf(clientSessionsKey, clientID), fmt.Sprintf(sessionIDKey, session.Page.ID, session.ID))
}
