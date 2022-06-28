package db

import (
	"context"
	"subscription/domain"
)

type DB interface {
	GetProduct(ctx context.Context, id string) ([]domain.Product, error)
	Disconnect(ctx context.Context) error
}
