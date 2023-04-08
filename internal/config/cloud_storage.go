package config

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
)

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


// ORIGINAL METHOD BEFORE REFACTORING
// DELETE THIS METHOD AFTER REFACTORING

// InitStorage initializes a Google Cloud Storage client.
// func InitStorage(ctx context.Context) (*storage.Client, error) {
// 	// TODO: Replace with your own Google Cloud Storage project credentials file path
// 	opt := option.WithCredentialsFile("<path/to/your/credentials/file>")
// 	client, err := storage.NewClient(ctx, opt)
// 	if err != nil {
// 		return nil, fmt.Errorf("error initializing Google Cloud Storage client: %v", err)
// 	}
// 	return client, nil
// }
