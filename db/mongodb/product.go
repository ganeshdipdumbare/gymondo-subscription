package mongodb

import (
	"context"

	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represent mongodb record from Product collection
type Product struct {
	Id                 primitive.ObjectID `bson:"_id"`
	Name               string             `bson:"name"`
	SubscriptionPeriod uint               `bson:"subscription_period"`
	Price              float64            `bson:"price"`
	TaxPercentage      float64            `bson:"tax_percentage"`
}

// createMongoProductRecord creates db product record from domain product
func createMongoProductRecord(p *domain.Product) *Product {
	return &Product{}
}

// createDomainProductRecord creates domain product record from db product
func createDomainProductRecord(p *Product) *domain.Product {
	return &domain.Product{}
}

func (m *mongoDetails) GetProduct(ctx context.Context, id string) ([]domain.Product, error) {
	return nil, nil
}
