package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRateLimitingMiddleware(t *testing.T) {
	// Create a mock handler to serve as the "next" handler in the middleware chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a new request with a URL that matches the rate-limited endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)

	// Call the middleware function with the mock handler and the new request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RateLimitingMiddleware(mockHandler))
	handler.ServeHTTP(rr, req)

	// Assert that the middleware function returns the expected response status code
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestRateLimitingMiddleware_RateLimitExceeded(t *testing.T) {
	// Create a mock handler to serve as the "next" handler in the middleware chain
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a new request with a URL that matches the rate-limited endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)

	// Create a rate limiter with a usage that exceeds the limit
	rateLimiter := &RateLimiter{
		Usage: 101,
		Limit: 100,
		Last:  time.Now().Add(-1 * time.Minute),
	}

	// Use a sync.Map to store the rate limiter information for the endpoint
	endpointLimiters := new(sync.Map)
	endpointLimiters.Store("/example", rateLimiter)

	// Call the middleware function with the mock handler and the new request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RateLimitingMiddleware(mockHandler))
	handler.ServeHTTP(rr, req)

	// Assert that the middleware function returns the expected response status code
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusTooManyRequests)
	}
}
