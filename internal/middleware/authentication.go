package middleware

import (
	"encoding/json"
	"net/http"
)

// AuthenticationMiddleware is a middleware function that checks the request
// for valid authentication credentials. If the request is not authenticated,
// the middleware returns an error response.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validSource(r) {
			// If the request is not authenticated, return an error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)			
			json.NewEncoder(w).Encode(`{"status": "error", "message": "Unauthorized request. Please verify you are making a request through a verified channels (i.e RapidAPI, Postman API Marketplace, etc..)."}`)
			return
		}

		// Check the request for valid authentication credentials
		if !validAuthentication(r) {
			// If the request is not authenticated, return an error response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(`{"status": "error", "message": "Unauthorized request. Please provide valid authentication credentials,"}`)
			return
		}

		// If the request is authenticated, call the next middleware/handler in the chain
		next.ServeHTTP(w, r)
	})
}

// validSource checks the request for valid source credentials. 
// If the request is not authenticated, the function returns false.
func validSource(r *http.Request) bool {
	// Implement your source logic here
	// Example: check the "Authorization" header for a valid source/referrer
	authHeader := r.Header.Get("X-RapidAPI-Proxy-Secret")
	if authHeader == "" {
		return false
	}

	return true
}

// validAuthentication checks the request for valid authentication credentials.
// If the request is not authenticated, the function returns false.
func validAuthentication(r *http.Request) bool {
	// Implement authentication logic here

	// Validate the token and return true if it's valid, false otherwise
	// Example: use a JWT library to decode and validate the token
	return true
}
