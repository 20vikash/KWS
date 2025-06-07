package main

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	// This map holds rate limiters per IP address.
	visitors  = make(map[string]*rate.Limiter)
	mu        sync.Mutex
	rateLimit = rate.Every(10 * time.Second) // 1 token every 10 seconds
	burst     = 5                            // bucket size
)

// Get the limiter for the IP address or create one
func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burst)
		visitors[ip] = limiter
	}
	return limiter
}
