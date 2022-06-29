package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserSubscription represent mongodb record from user_subscription collection
type UserSubscription struct {
	Id             primitive.ObjectID `bson:"_id"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      *time.Time         `bson:"updated_at,omitempty"`
	Email          string             `bson:"email"`
	ProductName    string             `bson:"product_name"`
	StartDate      time.Time          `bson:"start_date"`
	EndDate        time.Time          `bson:"end_date"`
	Price          float64            `bson:"price"`
	Tax            float64            `bson:"tax"`
	Status         string             `bson:"status"`
	PauseStartDate *time.Time         `bson:"pause_start_date,omitempty"`
}

// createDomainProductRecord creates domain product record from db product
func createDBUserSubscriptionRecord(us *domain.UserSubscription) (*UserSubscription, error) {
	if us == nil {
		return nil, fmt.Errorf("invalid input")
	}

	idHex, err := primitive.ObjectIDFromHex(us.ID)
	if err != nil {
		return nil, err
	}

	userSubscription := &UserSubscription{
		Id:          idHex,
		CreatedAt:   us.CreatedAt,
		Email:       us.Email,
		ProductName: us.ProductName,
		StartDate:   us.StartDate,
		EndDate:     us.EndDate,
		Price:       us.Price,
		Tax:         us.Tax,
		Status:      string(us.Status),
	}

	if us.UpdatedAt != nil {
		userSubscription.UpdatedAt = us.UpdatedAt
	}

	if us.PauseStartDate != nil {
		userSubscription.PauseStartDate = us.PauseStartDate
	}
	return userSubscription, nil
}

// SaveSubscription create new subscription if not present in the database otherwise update and return the subscription object
func (m *mongoDetails) SaveSubscription(ctx context.Context, us *domain.UserSubscription) (*domain.UserSubscription, error) {
	userSubscription, err := createDBUserSubscriptionRecord(us)
	if err != nil {
		return nil, err
	}
	opts := options.Replace().SetUpsert(true)
	filter := primitive.M{
		"_id": userSubscription.Id,
	}

	_, err = m.UserSubscriptionCollection.ReplaceOne(ctx, filter, userSubscription, opts)
	if err != nil {
		return nil, err
	}
	return us, nil
}
