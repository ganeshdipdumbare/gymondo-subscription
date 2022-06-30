package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/db"
	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
)

var (
	NilArgErr          = errors.New("nil value not allowed")
	InvalidArgErr      = errors.New("invalid argument")
	NotFoundErr        = errors.New("not found")
	NotAllowedArgErr   = errors.New("not allowed")
	StatusUnchangedErr = errors.New("status is unchanged")
)

//go:generate mockgen -destination=../mocks/mock_app.go -package=mocks github.com/ganeshdipdumbare/gymondo-subscription/app App
// App interface which consists of business logic/use cases
type App interface {
	GetProduct(ctx context.Context, id string) ([]domain.Product, error)
	BuySubscription(ctx context.Context, productID string, emailID string) (*domain.UserSubscription, error)
	GetSubscriptionByID(ctx context.Context, id string) (*domain.UserSubscription, error)
	UpdateSubscriptionStatusByID(ctx context.Context, id string, status domain.SubscriptionStatus) (*domain.UserSubscription, error)
}

type appDetails struct {
	database db.DB
}

// NewApp creates new app instance
func NewApp(database db.DB) (App, error) {
	if database == nil {
		return nil, fmt.Errorf("database %w", NilArgErr)
	}

	return &appDetails{
		database: database,
	}, nil
}

// GetProduct fetches product for given id, if id is not given then returns all the products
// returns invalid argument error if db returns invalid arg error
func (a *appDetails) GetProduct(ctx context.Context, id string) ([]domain.Product, error) {
	records, err := a.database.GetProduct(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, db.InvalidArgErr):
			return nil, fmt.Errorf("get product failed:%s %w", err.Error(), InvalidArgErr)
		default:
			return nil, err
		}
	}

	return records, nil
}

// BuySubscription subscription for given user id will be created for given product id
// returns invalid argument error if productID or emailID is empty
func (a *appDetails) BuySubscription(ctx context.Context, productID string, emailID string) (*domain.UserSubscription, error) {
	if productID == "" || emailID == "" {
		return nil, InvalidArgErr
	}

	records, err := a.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("%v %w", productID, NotFoundErr)
	}

	product := records[0]
	timeNow := time.Now().UTC()
	userSubscription := &domain.UserSubscription{
		CreatedAt:   timeNow,
		Email:       emailID,
		ProductName: product.Name,
		StartDate:   timeNow,
		EndDate:     timeNow.AddDate(0, int(product.SubscriptionPeriod), 0),
		Price:       product.Price,
		Status:      domain.SubscriptionStatusActive,
		Tax:         product.Price * product.TaxPercentage / 100,
	}

	return a.database.SaveSubscription(ctx, userSubscription)
}

// GetSubscriptionByID return subscription for given subscription id
// returns invalid argument if id is empty
func (a *appDetails) GetSubscriptionByID(ctx context.Context, id string) (*domain.UserSubscription, error) {
	if id == "" {
		return nil, InvalidArgErr
	}

	subscriptionDetails, err := a.database.GetSubscriptionByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, db.InvalidArgErr):
			return nil, fmt.Errorf("invalid argument:%s %w", err.Error(), InvalidArgErr)
		default:
			return nil, err
		}
	}
	return subscriptionDetails, nil
}

// UpdateSubscriptionStatusByID update subscription status by id
// status can be changed from active to cancelled or paused
// paused subscription can be unpaused/active or cancelled
// cancelled subscription status cannot be changed
func (a *appDetails) UpdateSubscriptionStatusByID(ctx context.Context, id string, status domain.SubscriptionStatus) (*domain.UserSubscription, error) {
	if id == "" {
		return nil, InvalidArgErr
	}

	subscriptionDetails, err := a.database.GetSubscriptionByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, db.InvalidArgErr):
			return nil, fmt.Errorf("invalid argument:%s %w", err.Error(), InvalidArgErr)
		default:
			return nil, err
		}
	}

	// check if the status is being changed
	if status == subscriptionDetails.Status {
		return nil, StatusUnchangedErr
	}

	timeNow := time.Now().UTC()
	updatedSubscriptionDetails := *subscriptionDetails
	updatedSubscriptionDetails.Status = status
	updatedSubscriptionDetails.UpdatedAt = &timeNow

	// check subscription's current status
	switch subscriptionDetails.Status {
	case domain.SubscriptionStatusCancelled:
		return nil, fmt.Errorf("cancelled subscription status change %w", NotAllowedArgErr)
	case domain.SubscriptionStatusPaused:
		if status == domain.SubscriptionStatusActive {
			pausedPeriod := timeNow.Sub(*subscriptionDetails.PauseStartDate)
			updatedSubscriptionDetails.EndDate = subscriptionDetails.EndDate.Add(pausedPeriod)
		}
	case domain.SubscriptionStatusActive:
		if status == domain.SubscriptionStatusPaused {
			updatedSubscriptionDetails.PauseStartDate = &timeNow
		}
	}

	return a.database.SaveSubscription(ctx, &updatedSubscriptionDetails)
}
