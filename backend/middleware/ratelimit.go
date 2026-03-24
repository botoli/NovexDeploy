package middleware

import (
	"net/http"
	"sync"
	"time"
)

type limiterEntry struct {
	count int
	reset time.Time
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	var mu sync.Mutex
	entries := make(map[string]*limiterEntry)
	limit := 120
	window := time.Minute

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		now := time.Now()

		mu.Lock()
		entry, ok := entries[key]
		if !ok || now.After(entry.reset) {
			entry = &limiterEntry{count: 0, reset: now.Add(window)}
			entries[key] = entry
		}
		entry.count++
		count := entry.count
		mu.Unlock()

		if count > limit {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
