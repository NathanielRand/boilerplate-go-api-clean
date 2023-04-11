package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AuthenticationMiddleware is a middleware function that checks the request
// for valid authentication credentials. If the request is not authenticated,
// the middleware returns an error response.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if !validSources(r) {
		// 	// If the request is not authenticated, return an error response
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	json.NewEncoder(w).Encode(`{"status": "error", "message": "Unauthorized request. Please verify you are making a request through a verified channels (i.e RapidAPI, Postman API Marketplace, etc..)."}`)
		// 	return
		// }

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

// validSource checks the request for valid source/referrer credentials.
func validSource(r *http.Request) bool {
	// Check if the request is coming from a valid source
	// via the request headers.
	// Example: X-RapidAPI-Proxy-Secret

	// Replace comparison value with environment variable
	if r.Header.Get("X-RapidAPI-Proxy-Secret") == "78f5b3e0-d3d0-11ed-bf92-43930995aeef" {
		return true
	} else {
		return false
	}
}

// validSources checks the request for valid sources/referrers credentials.
func validSources(r *http.Request) bool {
	// Define a map of valid sources
	validSources := map[string]string{
		"X-RapidAPI-Proxy-Secret": "78f5b3e0-d3d0-11ed-bf92-43930995aeef",
	}

	// Check the request for valid source credentials and return true if it's valid, false otherwise
	// Loop through the request headers and check for a valid source/referrer
	validSourceFound := false
	for key := range r.Header {
		value := r.Header.Get(key) // Get the value of the header
		fmt.Println("Checking request header: " + key + " = " + value + " ")
		// Check if the request header matches a valid source
		if validSourceValue, ok := validSources[key]; ok {
			// Check if the request header value matches a valid source
			if strings.EqualFold(value, validSourceValue) {
				fmt.Println("Valid source found: " + key + " = " + value)
				validSourceFound = true
			}
		}
	}

	// Return true if a valid source is found, false otherwise
	return validSourceFound
}

// validAuthentication checks the request for valid authentication credentials.
// If the request is not authenticated, the function returns false.
func validAuthentication(r *http.Request) bool {
	// Implement authentication logic here

	// Validate the token and return true if it's valid, false otherwise
	// Example: use a JWT library to decode and validate the token
	return true
}
