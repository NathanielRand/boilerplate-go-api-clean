package main

import (
	"log"

	"github.com/NathanielRand/morebytes-templates/go/api/clean-layered/cmd/server"
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
