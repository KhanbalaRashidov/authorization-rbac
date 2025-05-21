package cache

import (
"log"
"time"
)

func StartTokenCleanupService(repo *TokenRepo, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				repo.CleanupExpired()
				log.Println("🧹 Token cleanup completed")
			}
		}
	}()
}