package config

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
)

// Initialize the Google Cloud Storage client
func init() {
	// Initialize the Google Cloud Storage client
	var err error
	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}
}

// GetStorageClient returns the global storage client instance
func GetStorageClient() *storage.Client {
	return storageClient
}
