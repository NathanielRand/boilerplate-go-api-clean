package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware is a middleware function that logs the
// request details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		headersJSON, _ := json.Marshal(r.Header)
		logMessage := fmt.Sprintf("[%s] %s %s %s | User-Agent: %s | Headers: %s",
			time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI, r.RemoteAddr,
			r.Header.Get("User-Agent"), headersJSON)
		fmt.Println(logMessage)

		next.ServeHTTP(w, r)
	})
}

