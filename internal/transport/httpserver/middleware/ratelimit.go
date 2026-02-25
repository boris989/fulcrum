package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	rps     int
	burst   int
}

func NewRateLimiter(rps int, burst int) *rateLimiter {
	rl := &rateLimiter{
		clients: make(map[string]*client),
		rps:     rps,
		burst:   burst,
	}

	go rl.cleanup()

	return rl
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > 5*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		rl.mu.Lock()

		c, exists := rl.clients[ip]
		if !exists {
			c = &client{
				limiter:  rate.NewLimiter(rate.Limit(rl.rps), rl.burst),
				lastSeen: time.Now(),
			}
			rl.clients[ip] = c
		}
		c.lastSeen = time.Now()
		rl.mu.Unlock()

		if !c.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
