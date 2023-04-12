package main

import (
	"log"

	"github.com/NathanielRand/webchest-image-converter-api/cmd/server"
)

func main() {
	// Log the start of the server
	log.Println("Starting server...")

	// Start the server
	err := server.StartServer()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Log the end of the server
	log.Println("Server stopped")
}
