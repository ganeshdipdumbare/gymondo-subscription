package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/ganeshdipdumbare/gymondo-subscription/db"
	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
)

var (
	NilArgErr     = errors.New("nil value not allowed")
	InvalidArgErr = errors.New("invalid argument")
)

//go:generate mockgen -destination=../mocks/mock_app.go -package=mocks github.com/ganeshdipdumbare/gymondo-subscription/app App
// App interface which consists of business logic/use cases
type App interface {
	GetProduct(ctx context.Context, id string) ([]domain.Product, error)
}

type appDetails struct {
	database db.DB
}

// NewApp creates new app instance
func NewApp(database db.DB) (App, error) {
	if database == nil {
		return nil, fmt.Errorf("invalid database: %w", NilArgErr)
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
			return nil, fmt.Errorf("invalid argument:%s %w", err.Error(), InvalidArgErr)
		default:
			return nil, err
		}
	}

	return records, nil
}
