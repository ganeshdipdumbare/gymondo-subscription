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

// createDomainProductRecord creates domain product record from db product
func createDomainProductRecord(p *Product) *domain.Product {
	return &domain.Product{
		ID:                 p.Id.Hex(),
		Name:               p.Name,
		SubscriptionPeriod: p.SubscriptionPeriod,
		Price:              p.Price,
		TaxPercentage:      p.TaxPercentage,
	}
}

// createDomainProductRecordSl create domain Product slice from input product slice
func createDomainProductRecordSl(p []Product) []domain.Product {
	products := []domain.Product{}
	for _, v := range p {
		product := &v
		products = append(products, *createDomainProductRecord(product))
	}
	return products
}

// GetProduct returns product for given id, if id is not given, will return all the products
func (m *mongoDetails) GetProduct(ctx context.Context, id string) ([]domain.Product, error) {
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	productRecords := []Product{}
	err = m.getAllDocuments(ctx, m.ProductCollection, primitive.M{"_id": idHex}, &productRecords)
	if err != nil {
		return nil, err
	}

	return createDomainProductRecordSl(productRecords), nil
}
