package server

import (
	"net/http"
	"time"

	// "github.com/NathanielRand/boilerplate-go-api-clean/internal/config"
	"github.com/NathanielRand/boilerplate-go-api-clean/internal/routes"
)

// Start
func StartServer() error {

	// Get the router from the routes package
	router := routes.SetupRouter()

	// Get the port from the environment variables
	port := ":8080"

	// Load the environment variables
	// err := config.Load() // Load environment variables from .env file
	// if err != nil {
	// 	// Log the error
	// 	log.Printf("Error loading environment variables: %v", err)
	// 	return err
	// }

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
