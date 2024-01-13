package repositories

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/db"
	"github.com/NathanielRand/boilerplate-go-api-clean/internal/models"
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

// // CreateUser creates a new user in Firestore.
// func (r *FirestoreRepository) CreateUser(ctx context.Context, userRealIP string) (*models.User, error) {
// 	// Create a new user
// 	user := models.NewUser(userRealIP, username, key, forwaredIP, forwaredHost, subscription, affiliation, email, platform, quota, rateLimit)

// 	// Save the user to Firestore
// 	ref := r.client.NewRef(fmt.Sprintf("api-image-converter-users/%s", user.ID))
// 	if err := ref.Set(ctx, user); err != nil {
// 		return nil, err
// 	}

// 	// Return the user
// 	return user, nil
// }

// GetUserByID retrieves a user by ID from Firestore.
func (r *FirestoreRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	ref := r.client.NewRef(fmt.Sprintf("api-image-converter-users/%s", userID))
	if err := ref.Get(ctx, &user); err != nil {
		return nil, err
	}
	user.ID = userID
	return &user, nil
}

// GetUserByEmail retrieves a user by email from Firestore.
func (r *FirestoreRepository) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	var user models.User
	ref := r.client.NewRef("api-image-converter-users")
	if err := ref.OrderByChild("email").EqualTo(userEmail).LimitToFirst(1).Get(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByAPIKey retrieves a user by API key from Firestore.
func (r *FirestoreRepository) GetUserByAPIKey(ctx context.Context, userAPIKey string) (*models.User, error) {
	var user models.User
	ref := r.client.NewRef("api-image-converter-users")
	if err := ref.OrderByChild("api_key").EqualTo(userAPIKey).LimitToFirst(1).Get(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByRealIP retrieves a user by real IP from Firestore.
func (r *FirestoreRepository) GetUserByRealIP(ctx context.Context, userRealIP string) (*models.User, error) {
	var user models.User
	ref := r.client.NewRef("api-image-converter-users")
	if err := ref.OrderByChild("real_ip").EqualTo(userRealIP).LimitToFirst(1).Get(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckUserRealIP checks the real IP of a user in Firestore.
func (r *FirestoreRepository) CheckUserRealIP(ctx context.Context, userRealIP string) (string, error) {
	user, err := r.GetUserByRealIP(ctx, userRealIP)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// CheckUserPlan checks the subscription of a user in Firestore.
func (r *FirestoreRepository) CheckUserSubscription(ctx context.Context, userID string) (string, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.Subscription, nil
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
// func (r *FirestoreRepository) CheckUserReferrer(ctx context.Context, userID string) (string, error) {
// 	user, err := r.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return user.Referrer, nil
// }

// UpdateUser updates a user in Firestore.
func (r *FirestoreRepository) UpdateUser(ctx context.Context, userID string, user *models.User) error {
	ref := r.client.NewRef(fmt.Sprintf("api-image-converter-users/%s", userID))
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

// UpdateUserLoyalty updates the loyalty score status of a user in Firestore.
func (r *FirestoreRepository) UpdateUserLoyaltyScore(ctx context.Context, userID string, loyaltyStatus string) error {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.LoyaltyScore = loyaltyStatus

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
