package middleware

import (
	"log"
	"net/http"

	"github.com/cinchprotocol/cinch-api/packages/core"
)

var eventIdCheckEndpoints = []string{
	"/api/v1/resource1",
}

func RequestEventId(cacheClient core.ICache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, endpoint := range eventIdCheckEndpoints {
				if r.URL.Path == endpoint {
					eventId := r.Header.Get("X-Event-ID")
					if eventId == "" {
						http.Error(w, "Missing X-Event-ID header", http.StatusBadRequest)
						return
					}
					keyExisted, err := cacheClient.Exists(eventId)
					if err != nil {
						log.Printf("Error checking Redis key existence: %s", err)
						break
					}
					if keyExisted == true {
						http.Error(w, "X-Event-ID already processed", http.StatusConflict)
						return
					}
					_, _ = cacheClient.Set(eventId, "1")
					break
				}
			}
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}
