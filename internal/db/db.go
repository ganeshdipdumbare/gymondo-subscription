package db

import (
	"context"
	"errors"

	"github.com/ganeshdipdumbare/gymondo-subscription/internal/domain"
)

var (
	InvalidArgErr     = errors.New("invalid argument")
	EmptyArgErr       = errors.New("empty argument not allowed")
	RecordNotFoundErr = errors.New("record not found")
)

// DB interface to interact with database
//
//go:generate mockgen -destination=../mocks/mock_db.go -package=mocks github.com/ganeshdipdumbare/gymondo-subscription/internal/db DB
type DB interface {
	GetProduct(ctx context.Context, id string) ([]domain.Product, error)
	SaveSubscription(ctx context.Context, subsciption *domain.UserSubscription) (*domain.UserSubscription, error)
	GetSubscriptionByID(ctx context.Context, id string) (*domain.UserSubscription, error)
	Disconnect(ctx context.Context) error
}
