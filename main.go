package main

import (
	"log"

	"github.com/NathanielRand/webchest-image-converter-api/cmd/server"
	"github.com/NathanielRand/webchest-image-converter-api/internal/config"
)

func main() {
	// Load the environment variables
	err := config.Load() // Load environment variables from .env file
	if err != nil {
		log.Fatal(err)
	}

	// Log the start of the server
	log.Println("Starting server...")

	// Start the server
	err = server.StartServer()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Log the end of the server
	log.Println("Server stopped")
}
