package repositories

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/db"
	"github.com/NathanielRand/webchest-image-converter-api/internal/models"
)

// FirestoreRepository is a repository that retrieves data from Firestore.
type FirestoreRepository struct {
	client *db.Client
}

// NewFirestoreRepository creates a new FirestoreRepository.
func NewFirestoreRepository(client *db.Client) *FirestoreRepository {
	return &FirestoreRepository{
		client: client,
	}
}

// CustomAction performs a custom action for a user in Firestore.
func (r *FirestoreRepository) CustomAction(ctx context.Context, userID string, actionType string) (string, error) {
	// Custom logic for performing a specific action based on actionType
	// For example, you can implement different actions based on the value of actionType

	// Example custom action:
	if actionType == "someAction" {
		// Perform some action
		return "Action completed", nil
	} else {
		return "", errors.New("Invalid action type")
	}
}

// GetUserByID retrieves a user by ID from Firestore.
func (r *FirestoreRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	ref := r.client.NewRef(fmt.Sprintf("users/%s", userID))
	if err := ref.Get(ctx, &user); err != nil {
		return nil, err
	}
	user.ID = userID
	return &user, nil
}

// CheckUserPlan checks the plan of a user in Firestore.
func (r *FirestoreRepository) CheckUserPlan(ctx context.Context, userID string) (string, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.Plan, nil
}

// CheckUserQuota checks the quota of a user in Firestore.
func (r *FirestoreRepository) CheckUserQuota(ctx context.Context, userID string) (int, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return user.Quota, nil
}

// CheckUserRateLimit checks the rate limit of a user in Firestore.
func (r *FirestoreRepository) CheckUserRateLimit(ctx context.Context, userID string) (int, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return user.RateLimit, nil
}

// CheckUserReferrer checks the referrer of a user in Firestore.
func (r *FirestoreRepository) CheckUserReferrer(ctx context.Context, userID string) (string, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.Referrer, nil
}

// UpdateUser updates a user in Firestore.
func (r *FirestoreRepository) UpdateUser(ctx context.Context, userID string, user *models.User) error {
	ref := r.client.NewRef(fmt.Sprintf("users/%s", userID))
	if err := ref.Set(ctx, user); err != nil {
		return err
	}
	return nil
}

// AddQuota adds quota to a user in Firestore.
func (r *FirestoreRepository) AddQuota(ctx context.Context, userID string, quota int) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Quota += quota

	err = r.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserKeys updates the keys of a user in Firestore.
func (r *FirestoreRepository) UpdateUserKeys(ctx context.Context, userID string, keys []string) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Keys = keys

	err = r.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserSpend updates the spend of a user in Firestore.
func (r *FirestoreRepository) UpdateUserSpend(ctx context.Context, userID string, spend float64) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Spend += spend

	err = r.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserLoyalty updates the loyalty status of a user in Firestore.
func (r *FirestoreRepository) UpdateUserLoyalty(ctx context.Context, userID string, loyaltyStatus string) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Loyalty = loyaltyStatus

	err = r.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserAffiliations updates the affiliations of a user in Firestore.
func (r *FirestoreRepository) UpdateUserAffiliations(ctx context.Context, userID string, affiliations []string) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Affiliations = affiliations

	err = r.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}
