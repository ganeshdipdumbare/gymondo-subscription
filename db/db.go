package db

import (
	"context"
	"errors"

	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
)

var (
	InvalidArgErr = errors.New("invalid argument")
	EmptyArgErr   = errors.New("empty argument not allowed")
	NotFoundErr   = errors.New("record not found")
)

//go:generate mockgen -destination=../mocks/mock_db.go -package=mocks github.com/ganeshdipdumbare/gymondo-subscription/db DB
// DB interface to interact with database
type DB interface {
	GetProduct(ctx context.Context, id string) ([]domain.Product, error)
	SaveSubscription(ctx context.Context, subsciption *domain.UserSubscription) (*domain.UserSubscription, error)
	GetSubscriptionByID(ctx context.Context, id string) (*domain.UserSubscription, error)
	Disconnect(ctx context.Context) error
}
