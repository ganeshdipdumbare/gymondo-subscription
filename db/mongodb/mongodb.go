package mongodb

import (
	"context"
	"fmt"
	"subscription/db"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	emptyArgErr = "empty %v not allowed"
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
		return nil, fmt.Errorf(emptyArgErr, "uri")
	}

	if dbName == "" {
		return nil, fmt.Errorf(emptyArgErr, "dbName")
	}

	client, err := connect(uri)
	if err != nil {
		return nil, err
	}

	productCollection := client.Database(dbName).Collection("product")
	userSubscriptionCollection := client.Database(dbName).Collection("user_subscription")

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

// Disconnect disconnects db connection using client, otherwise returns error
func (m *mongoDetails) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
