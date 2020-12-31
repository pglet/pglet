package page

import (
	"context"
	"log"
	"time"

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
		oldPages := store.GetLastUpdatedPages(time.Now().Add(-1 * time.Minute).Unix())
		if len(oldPages) > 0 {
			log.Println("OLD PAGES:", oldPages)
		}
	}
}
