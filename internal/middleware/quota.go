package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// UserQuota stores the usage and limit information for a user
type UserQuota struct {
	Usage  int64
	Limit  int64
	Last   time.Time
	Tokens int64
}

// QuotaMiddleware is a middleware function that enforces a quota
// for each user based on a token bucket algorithm.
func QuotaMiddleware(next http.Handler) http.Handler {
	// Your quota middleware logic goes here
	// Use a sync.Map to store the UserQuota information for each user
	userQuotas := new(sync.Map)

	// Return the middleware function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user's ID from the request
		userID := "someUserID" // Replace this with your code to extract the user ID from the request

		// Get or create the user's quota information
		userQuota, _ := userQuotas.LoadOrStore(userID, &UserQuota{
			Usage:  0,
			Limit:  100000000, // Replace with your desired limit
			Last:   time.Now(),
			Tokens: 5, // Replace with your desired initial tokens
		})

		// Calculate the time elapsed since the last request
		elapsed := time.Since(userQuota.(*UserQuota).Last)

		// Calculate the number of tokens that should be added to the bucket
		addTokens := int64(elapsed.Seconds() * float64(userQuota.(*UserQuota).Limit) / 3600)

		// Add the tokens to the bucket, up to the maximum limit
		userQuota.(*UserQuota).Tokens = min(userQuota.(*UserQuota).Tokens+addTokens, userQuota.(*UserQuota).Limit)

		// Check if the request can be served
		if userQuota.(*UserQuota).Tokens > 0 {
			// Serve the request
			userQuota.(*UserQuota).Tokens--
			userQuota.(*UserQuota).Usage++
			userQuota.(*UserQuota).Last = time.Now()

			// Print the quota information
			fmt.Printf("User %s: %d/%d tokens, %d/%d requests\n", userID, userQuota.(*UserQuota).Tokens, userQuota.(*UserQuota).Limit, userQuota.(*UserQuota).Usage, userQuota.(*UserQuota).Limit)

			next.ServeHTTP(w, r)
			return
		}

		// Return an error response if the quota has been exceeded
		http.Error(w, "quota exceeded", http.StatusTooManyRequests)
	})

}

// min returns the smaller of two integers
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
