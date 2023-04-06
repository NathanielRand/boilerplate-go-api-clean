package middleware

import (
	"net/http"
)

// AuthorizationMiddleware is a middleware function that checks the request
// for valid authorization credentials. If the request is not authorized,
// the middleware returns an error response.
func AuthorizationMiddleware(next http.Handler, requiredRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request for valid authorization credentials
		if !validAuthorization(r, requiredRole) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// If the request is authorized, call the next middleware/handler in the chain
		next.ServeHTTP(w, r)
	})
}

// validAuthorization checks the request for valid authorization credentials.
// If the request is not authorized, the function returns false.
func validAuthorization(r *http.Request, requiredRole string) bool {
	// Implement your authorization logic here
	// Example: check the authenticated user's role against the required role
	userRole := getUserRole(r)
	if userRole != requiredRole {
		return false
	}

	// If the user has the required role, return true
	return true
}

// getUserRole retrieves the authenticated user's role from the request context
// or from the authentication credentials (e.g. JWT claims).
func getUserRole(r *http.Request) string {
	// Implement your user role retrieval logic here
	// Example: retrieve the "role" claim from the JWT token
	return "admin" // replace with actual user role
}
