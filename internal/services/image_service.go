package services

import (
	"context"

	"github.com/NathanielRand/boilerplate-go-api-clean/internal/models"
	"github.com/NathanielRand/boilerplate-go-api-clean/internal/repositories"
)

// ImageService is a service that retrieves data from Firestore/Cloud Storage.
type ImageService struct {
	repo *repositories.FirestoreRepository
}

// NewImageService creates a new ImageService.
func NewImageService(repo *repositories.FirestoreRepository) *ImageService {
	return &ImageService{
		repo: repo,
	}
}

// GetImages retrieves all images from Firestore.
func (s *ImageService) GetImages(ctx context.Context) ([]*models.Image, error) {
	return s.repo.GetImages(ctx)
}

// AddImage adds an image to Firestore.
func (s *ImageService) AddImage(ctx context.Context, image *models.Image) error {
	return s.repo.AddImage(ctx, image)
}
