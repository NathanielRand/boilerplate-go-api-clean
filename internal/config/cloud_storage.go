package config

import (
	"context"
	"cloud.google.com/go/storage"
	"fmt"
	"google.golang.org/api/option"
)

func InitStorage(ctx context.Context) (*storage.Client, error) {
	// TODO: Replace with your own Google Cloud Storage project credentials file path
	opt := option.WithCredentialsFile("<path/to/your/credentials/file>")
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Google Cloud Storage client: %v", err)
	}
	return client, nil
}
