package repositories

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

// CloudStorageRepository is a repository for uploading images to Google Cloud Storage.
type CloudStorageRepository struct {
	bucket *storage.BucketHandle
}

// NewCloudStorageRepository creates a new CloudStorageRepository.
func NewCloudStorageRepository(bucketName string, client *storage.Client) (*CloudStorageRepository, error) {
	// Create a new bucket handle
	bucket := client.Bucket(bucketName)

	// Get the current time
	now := time.Now()

	// Calculate the time 6 minutes from now
	expirationTime := now.Add(6 * time.Minute)

	fmt.Printf("Expiration time: ", expirationTime)

	// Define the lifecycle rule
	rule := &storage.Lifecycle{
		Rules: []storage.LifecycleRule{
			{
				Action: storage.LifecycleAction{
					Type: "Delete",
				},
				Condition: storage.LifecycleCondition{
					CustomTimeBefore: expirationTime, // Custom time before
				},
			},
		},
	}

	// Update the bucket attributes
	attrs := storage.BucketAttrsToUpdate{
		Lifecycle: rule,
	}
	if _, err := bucket.Update(context.Background(), attrs); err != nil {
		return nil, err
	}

	// Return a new CloudStorageRepository
	return &CloudStorageRepository{
		bucket: bucket,
	}, nil
}


// UploadImage uploads an image to Google Cloud Storage.
func (r *CloudStorageRepository) UploadImage(ctx context.Context, name string, data io.Reader) (string, error) {
	// Create a new object in the bucket
	obj := r.bucket.Object(name)

	// Create a new writer for the object
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, data); err != nil {
		return "", err
	}

	// Close the writer
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Set the cache control header on the uploaded object
	attrs := storage.ObjectAttrsToUpdate{
		CacheControl: "public, max-age=31536000",
		CustomTime:   time.Now(),
	}

	// Update the object attributes
	if _, err := obj.Update(ctx, attrs); err != nil {
		return "", err
	}

	// Get the public URL of the uploaded object and return it
	url := strings.Replace(obj.ObjectName(), "/", "%2F", -1)
	return "https://storage.googleapis.com/" + obj.BucketName() + "/" + url, nil
}
