package server

import (
	"log"
	"net/http"
	"time"

	"github.com/NathanielRand/webchest-image-converter-api/internal/routes"
	"github.com/NathanielRand/webchest-image-converter-api/internal/config"
)

// Start
func StartServer() error {
	// Load the environment variables
	err := config.Load() // Load environment variables from .env file
	if err != nil {
		log.Fatal(err)
	}
	
	// Get the router from the routes package
	router := routes.SetupRouter()

	// Create an HTTP server with timeouts
	server := &http.Server{
		Addr:         config.Get("PORT"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the HTTP server
	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
