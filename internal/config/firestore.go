package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	// "google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
)

func init() {
	// Initialize the Google Cloud Firestore client
	ctx := context.Background()

	// Set your Google Cloud Platform project ID
	projectID := "webchest"

	// Create Firestore client with project ID
	// Additional Options (if needed due to deployment 
	// outside of google cloud): option.WithCredentialsFile("path/to/credentials.json")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Google Cloud Firestore client: %v", err)
	}
	firestoreClient = client
}

// GetFirestoreClient returns the global Firestore client instance
func GetFirestoreClient() *firestore.Client {
	return firestoreClient
}
