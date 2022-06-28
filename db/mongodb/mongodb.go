package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	productCollection          = "product"
	userSubscriptionCollection = "user_subscription"
)

type mongoDetails struct {
	client                     *mongo.Client
	dbName                     string
	ProductCollection          *mongo.Collection
	UserSubscriptionCollection *mongo.Collection
}

// NewMongoDB created new mongo db instance, returns error if input is invalid
func NewMongoDB(uri string, dbName string) (db.DB, error) {

	if uri == "" {
		return nil, fmt.Errorf("NewMongoDB: empty url %w", db.EmptyArgErr)
	}

	if dbName == "" {
		return nil, fmt.Errorf("NewMongoDB: empty db name %w", db.EmptyArgErr)
	}

	client, err := connect(uri)
	if err != nil {
		return nil, err
	}

	productCollection := client.Database(dbName).Collection(productCollection)
	userSubscriptionCollection := client.Database(dbName).Collection(userSubscriptionCollection)

	return &mongoDetails{
		client:                     client,
		dbName:                     dbName,
		ProductCollection:          productCollection,
		UserSubscriptionCollection: userSubscriptionCollection,
	}, nil
}

// connect connects to mongo db using client, returns error if fails
func connect(uri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getAllDocuments returns all the documents for matching filter from given collection, otherwise error
func (m *mongoDetails) getAllDocuments(ctx context.Context, collection *mongo.Collection, filter primitive.M, records interface{}) error {
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}

	return cur.All(ctx, records)
}

// Disconnect disconnects db connection using client, otherwise returns error
func (m *mongoDetails) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
