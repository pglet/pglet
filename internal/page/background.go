package page

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/pubsub"
	"github.com/pglet/pglet/internal/store"
)

func RunBackgroundTasks(ctx context.Context) {
	log.Println("Starting background tasks...")
	go CleanupPagesAndSessions()
	go CleanupExpiredClients()
}

func CleanupPagesAndSessions() {
	log.Println("Start background task to cleanup old pages and sessions")

	ticker := time.NewTicker(10 * time.Second)
	for {
		<-ticker.C

		sessions := store.GetExpiredSessions()
		if len(sessions) > 0 {
			log.Debugln("Deleting old sessions:", len(sessions))
			for _, fullSessionID := range sessions {
				pageID, sessionID := model.ParseSessionID(fullSessionID)

				page := store.GetPageByID(pageID)
				if page == nil {
					continue
				}

				// notify host client about expired session
				msg := NewMessageData("", PageEventToHostAction, &PageEventPayload{
					PageName:    page.Name,
					SessionID:   sessionID,
					EventTarget: "page",
					EventName:   "close",
				})

				for _, clientID := range store.GetSessionHostClients(pageID, sessionID) {
					pubsub.Send(clientChannelName(clientID), msg)
				}

				store.DeleteSession(pageID, sessionID)

				// delete page if no more sessions
				if !page.IsApp && len(store.GetPageSessions(pageID)) == 0 {
					store.DeletePage(pageID)
				}
			}
		}
	}
}

func CleanupExpiredClients() {
	log.Println("Start background task to cleanup expired clients")

	ticker := time.NewTicker(20 * time.Second)
	for {
		<-ticker.C

		clients := store.GetExpiredClients()
		for _, clientID := range clients {
			log.Debugln("Delete expired client:", clientID)
			store.DeleteExpiredClient(clientID)
		}
	}
}
