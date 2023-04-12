package server

import (
	"net/http"
	"time"

	"github.com/NathanielRand/webchest-image-converter-api/internal/routes"
)

const (
	port = ":8080"
)

func StartServer() error {
	// Get the router from the routes package
	router := routes.SetupRouter()

	// Create an HTTP server with timeouts
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the HTTP server
	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
