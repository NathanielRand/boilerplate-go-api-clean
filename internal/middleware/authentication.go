package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

// AuthenticationMiddleware is a middleware function that checks the request
// for valid authentication credentials. If the request is not authenticated,
// the middleware returns an error response.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request for valid source/referrer credentials
		// if !validSource(r) {
		// 	// If the request is not from a valid source, return an error response
		// 	w.Header().Set("Content-Type", "application/json")
		// 	json.NewEncoder(w).Encode(`{"status": "error", "message": "Unauthorized request. Please verify you are making a request through a verified channel (i.e RapidAPI, Postman API Marketplace, etc..)."}`)
		// 	return
		// }

		// Check the request for valid authentication credentials
		if !validAuthentication(r) {
			// If the request is not authenticated, return an error response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(`{"status": "error", "message": "Unauthorized request. Please provide valid authentication credentials."}`)
			return
		}

		// Check if the user exists in the database, and if the user is active
		// otherwise, create a new user in the database and continue,
		// if !userExists(r) {
		// 	// Get the user API key from the request headers
		// 	userRealIP := r.Header.Get("X-RapidAPI-Real-IP")

		// 	// Get a Firestore client
		// 	firestoreClient, err := config.GetFirestoreClient()

		// 	// if the user does not exist, create a new user in the database
		// 	// and continue
		// 	user, err := repositories.NewFirestoreRepository(firestoreClient).CreateUser(r.Context(), userRealIP)
		// 	if err != nil {
		// 		// If there is an error creating the user, return an error response
		// 		w.Header().Set("Content-Type", "application/json")
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		json.NewEncoder(w).Encode(`{"status": "error", "message": "Error creating user in database."}`)
		// 		return
		// 	}
		// }

		// If the request is authenticated, call the next middleware/handler in the chain
		next.ServeHTTP(w, r)
	})
}

func validSource(r *http.Request) bool {
	// Define a map of valid source keys and values
	validSources := map[string]string{
		"x-api-proxy-secret": "x8cc0e10-dd53-11ed-a321-315a0260571a",
		// Add more valid source keys and values here
	}

	// Iterate through the request headers
	for key, values := range r.Header {
		// Convert the key to lowercase
		keyLower := strings.ToLower(key)
		// Check if the key exists in the valid sources map
		if validValue, ok := validSources[keyLower]; ok {
			// Iterate through the values in the header
			for _, value := range values {
				// Convert the header value and valid value to lowercase
				valueLower := strings.ToLower(value)
				validValueLower := strings.ToLower(validValue)
				// Check if any of the values match the valid value
				if valueLower == validValueLower {
					return true
				}
			}
		}
	}
	return false
}

// validAuthentication checks the request for valid authentication credentials.
// If the request is not authenticated, the function returns false.
func validAuthentication(r *http.Request) bool {
	// Implement authentication logic here

	// Validate the token and return true if it's valid, false otherwise
	// Example: use a JWT library to decode and validate the token
	return true
}
