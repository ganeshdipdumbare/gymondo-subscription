package db

import (
	"context"

	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
)

//go:generate mockgen -destination=../mocks/mock_db.go -package=mocks github.com/ganeshdipdumbare/gymondo-subscription/db DB
// DB interface to interact with database
type DB interface {
	GetProduct(ctx context.Context, id string) ([]domain.Product, error)
	Disconnect(ctx context.Context) error
}
