package config

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

// InitFirestore initializes a Firestore client.
func InitFirestore(ctx context.Context) (*db.Client, error) {
	// TODO: Replace with your own Firebase project credentials file path
	opt := option.WithCredentialsFile("<path/to/your/credentials/file>")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firestore client: %v", err)
	}
	return client, nil
}
