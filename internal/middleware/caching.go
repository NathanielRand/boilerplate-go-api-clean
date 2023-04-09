package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// CachingMiddleware is a middleware function that checks the
// request against a caching policy. If the response is cached,
// the middleware returns the cached response. Otherwise, the
// middleware calls the next handler in the chain and caches
// the response.
func CachingMiddleware(next http.Handler) http.Handler {
	// Initialize the cache with a default expiration time of 5 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a cache key based on the request URL and method
		cacheKey := r.RemoteAddr + r.RequestURI + r.FormValue("image") + r.FormValue("format")

		// Check if the response is already cached
		if cachedResponse, found := c.Get(cacheKey); found {
			// If the response is cached, return the cached response
			cachedResponseBytes, ok := cachedResponse.([]byte)
			if !ok {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(cachedResponseBytes)
			return
		}

		// Otherwise, call the next handler in the chain and cache the response
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		// Copy the response body to a buffer to cache the response
		buffer := bytes.NewBuffer(nil)
		_, err := io.Copy(buffer, recorder.Body)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Cache the response for future requests
		c.Set(cacheKey, buffer.Bytes(), cache.DefaultExpiration)

		// Write the response to the original writer
		recorder.Result().Body.Close()
		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(recorder.Code)
		w.Write(buffer.Bytes())
	})
}
