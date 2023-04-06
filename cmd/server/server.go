package server

import (
	"net/http"
	"time"

	"github.com/NathanielRand/morebytes-templates/go/api/clean-layered/internal/routes"
)

func StartServer() error {
	// Get the router from the routes package
	router := routes.SetupRouter()

	// Create an HTTP server with timeouts
	server := &http.Server{
		Addr:         ":8080",
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
