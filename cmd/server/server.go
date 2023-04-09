package server

import (
	"net/http"
	"time"

	"github.com/NathanielRand/webchest-image-converter-api/internal/routes"
)

func StartServer() error {
	// Get the router from the routes package
	router := routes.SetupRouter()

	// Create an HTTP server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 100 * time.Second,
	}

	// Start the HTTP server
	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
