package page

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pglet/pglet/internal/model"
	"github.com/pglet/pglet/internal/store"
)

func RunBackgroundTasks(ctx context.Context) {
	log.Println("Starting background tasks...")
	go CleanupPagesAndSessions()
}

func CleanupPagesAndSessions() {
	log.Println("Start background task to cleanup old pages and sessions")

	ticker := time.NewTicker(10 * time.Second)
	for {
		<-ticker.C

		//log.Println("Cleanup pages and sessions!")
		sessions := store.GetExpiredSessions()
		if len(sessions) > 0 {
			log.Println("Deleting old sessions:", len(sessions))
			for _, fullSessionID := range sessions {
				pageID, sessionID := model.ParseSessionID(fullSessionID)
				store.DeleteSession(pageID, sessionID)

				// delete page if no more sessions
				if len(store.GetPageSessions(pageID)) == 0 {
					store.DeletePage(pageID)
				}
			}
		}
	}
}
