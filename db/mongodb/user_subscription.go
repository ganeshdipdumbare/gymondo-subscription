package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/db"
	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserSubscription represent mongodb record from user_subscription collection
type UserSubscription struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
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

// createDomainProductRecord creates db UserSbuscription record from domain record
func createDBUserSubscriptionRecord(us *domain.UserSubscription) (*UserSubscription, error) {
	if us == nil {
		return nil, db.InvalidArgErr
	}

	userSubscription := &UserSubscription{
		CreatedAt:   us.CreatedAt,
		Email:       us.Email,
		ProductName: us.ProductName,
		StartDate:   us.StartDate,
		EndDate:     us.EndDate,
		Price:       us.Price,
		Tax:         us.Tax,
		Status:      string(us.Status),
	}

	if us.ID != "" {
		idHex, err := primitive.ObjectIDFromHex(us.ID)
		if err != nil {
			return nil, db.InvalidArgErr
		}
		userSubscription.Id = idHex
	}

	if us.UpdatedAt != nil {
		userSubscription.UpdatedAt = us.UpdatedAt
	}

	if us.PauseStartDate != nil {
		userSubscription.PauseStartDate = us.PauseStartDate
	}
	return userSubscription, nil
}

// createDomainUserSubscriptionRecord creates domain UserSubscription record from db record
func createDomainUserSubscriptionRecord(us *UserSubscription) (*domain.UserSubscription, error) {
	if us == nil {
		return nil, db.InvalidArgErr
	}

	userSubscription := &domain.UserSubscription{
		ID:          us.Id.Hex(),
		CreatedAt:   us.CreatedAt,
		Email:       us.Email,
		ProductName: us.ProductName,
		StartDate:   us.StartDate,
		EndDate:     us.EndDate,
		Price:       us.Price,
		Tax:         us.Tax,
		Status:      domain.SubscriptionStatus(us.Status),
	}

	if us.UpdatedAt != nil {
		userSubscription.UpdatedAt = us.UpdatedAt
	}

	if us.PauseStartDate != nil {
		userSubscription.PauseStartDate = us.PauseStartDate
	}

	return userSubscription, nil
}

// SaveSubscription create new subscription if not present in the database otherwise update and return the subscription record
func (m *mongoDetails) SaveSubscription(ctx context.Context, us *domain.UserSubscription) (*domain.UserSubscription, error) {
	userSubscription, err := createDBUserSubscriptionRecord(us)
	if err != nil {
		return nil, err
	}

	var recordID primitive.ObjectID
	if !userSubscription.Id.IsZero() {
		recordID = userSubscription.Id
	} else {
		recordID = primitive.NewObjectID()
	}
	us.ID = recordID.Hex()
	filter := primitive.M{
		"_id": recordID,
	}

	opts := options.Replace().SetUpsert(true)
	_, err = m.UserSubscriptionCollection.ReplaceOne(ctx, filter, userSubscription, opts)
	if err != nil {
		return nil, err
	}

	return us, nil
}

// GetSubscriptionByID return subscription for given id
func (m *mongoDetails) GetSubscriptionByID(ctx context.Context, id string) (*domain.UserSubscription, error) {

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("id %w", db.InvalidArgErr)
	}

	filter := primitive.M{"_id": idHex}
	var record UserSubscription
	err = m.UserSubscriptionCollection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, db.RecordNotFoundErr
		}
		return nil, err
	}
	return createDomainUserSubscriptionRecord(&record)
}
