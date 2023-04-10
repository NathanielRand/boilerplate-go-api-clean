package middleware

import (
	"fmt"
	"net/http"
)

// AuthenticationMiddleware is a middleware function that checks the request
// for valid authentication credentials. If the request is not authenticated,
// the middleware returns an error response.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request for valid authentication credentials
		if !validAuthentication(r) {
			fmt.Fprintf(w, `{"status": "error", "message": "Unauthorized request. Please provide valid authentication credentials and/or verify you are making the request through the proper infrastructure (i.e RapidAPI, Postman API Marketplace, etc..)."}`)
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the request is authenticated, call the next middleware/handler in the chain
		next.ServeHTTP(w, r)
	})
}

// validAuthentication checks the request for valid authentication credentials.
// If the request is not authenticated, the function returns false.
func validAuthentication(r *http.Request) bool {
	// Implement your authentication logic here
	// Example: check the "Authorization" header for a valid token
	authHeader := r.Header.Get("X-Authorization")
	if authHeader == "" {
		return false
	}

	// Verify the request is coming from RapidAPI 
	// (replace with whitelist of IPs supporting the included infrastructure)
	RapidAPIHeader := r.Header.Get("X-RapidAPI-Proxy-Secret")
	if RapidAPIHeader == "" {
		return false
	}

	// Validate the token and return true if it's valid, false otherwise
	// Example: use a JWT library to decode and validate the token
	return true
}
