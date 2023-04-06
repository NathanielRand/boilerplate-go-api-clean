package repositories

import (
	"context"
	"io"
	"strings"

	"cloud.google.com/go/storage"
)

// CloudStorageRepository is a repository for uploading images to Google Cloud Storage.
type CloudStorageRepository struct {
	bucket *storage.BucketHandle
}

// NewCloudStorageRepository creates a new CloudStorageRepository.
func NewCloudStorageRepository(bucketName string, client *storage.Client) *CloudStorageRepository {
	return &CloudStorageRepository{
		bucket: client.Bucket(bucketName),
	}
}

// UploadImage uploads an image to Google Cloud Storage.
func (r *CloudStorageRepository) UploadImage(ctx context.Context, name string, data io.Reader) (string, error) {
	// Create a new object in the bucket and upload the image data
	obj := r.bucket.Object(name)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, data); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Set the cache control header on the uploaded object
	attrs := storage.ObjectAttrsToUpdate{
		CacheControl: "public, max-age=31536000",
	}
	if _, err := obj.Update(ctx, attrs); err != nil {
		return "", err
	}

	// Get the public URL of the uploaded object
	url := strings.Replace(obj.ObjectName(), "/", "%2F", -1)
	return "https://storage.googleapis.com/" + r.bucket.Name() + "/" + url, nil
}
