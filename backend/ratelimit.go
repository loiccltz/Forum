package backend

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visits map[string][]time.Time
	mu     sync.Mutex
}

var limiter = RateLimiter{
	visits: make(map[string][]time.Time),
}

func LimitRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()

		var recentRequests []time.Time
		for _, t := range limiter.visits[ip] {
			if now.Sub(t) < time.Minute {
				recentRequests = append(recentRequests, t)
			}
		}
		limiter.visits[ip] = recentRequests

		if len(limiter.visits[ip]) >= 100 {
			http.Error(w, "⛔ Trop de requêtes, veuillez patienter.", http.StatusTooManyRequests)
			return
		}

		limiter.visits[ip] = append(limiter.visits[ip], now)

		next.ServeHTTP(w, r)
	})
}
