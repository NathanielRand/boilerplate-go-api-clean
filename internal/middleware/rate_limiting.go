package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter stores the usage and limit information for a rate-limited endpoint
type RateLimiter struct {
	Usage int64
	Limit int64
	Last  time.Time
	Mutex sync.Mutex
}

// RateLimitingMiddleware is a middleware function that checks the
// request against a rate limit policy. If the request exceeds the
// rate limit, the middleware returns an error response.
func RateLimitingMiddleware(next http.Handler) http.Handler {
	// Your rate limiting middleware logic goes here
	// Use a sync.Map to store the RateLimiter information for each endpoint
	endpointLimiters := new(sync.Map)

	// Return the middleware function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the endpoint's URL from the request
		endpoint := r.URL.Path

		// Get or create the endpoint's rate limiter
		endpointLimiter, _ := endpointLimiters.LoadOrStore(endpoint, &RateLimiter{
			Usage: 0,
			Limit: 32, // Replace with new limit
			Last:  time.Now(),
		})

		// Acquire a mutex lock to avoid race conditions
		endpointLimiter.(*RateLimiter).Mutex.Lock()
		defer endpointLimiter.(*RateLimiter).Mutex.Unlock()

		// Calculate the time elapsed since the last request
		elapsed := time.Since(endpointLimiter.(*RateLimiter).Last)

		// Calculate the number of requests that should be allowed
		numRequests := int64(elapsed.Seconds() * float64(endpointLimiter.(*RateLimiter).Limit))

		// Update the rate limiter information
		endpointLimiter.(*RateLimiter).Usage = max(endpointLimiter.(*RateLimiter).Usage-numRequests, 0) + 1
		endpointLimiter.(*RateLimiter).Last = time.Now()

		// Check if the request can be served
		if endpointLimiter.(*RateLimiter).Usage <= endpointLimiter.(*RateLimiter).Limit {
			// Serve the request
			next.ServeHTTP(w, r)
			return
		}

		// Return an error response if the rate limit has been exceeded
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
	})

}

// max returns the larger of two integers
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
