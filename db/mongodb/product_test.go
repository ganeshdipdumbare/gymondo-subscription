package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_createDomainProductRecord(t *testing.T) {
	productIDHex := primitive.NewObjectID()
	type args struct {
		p *Product
	}
	tests := []struct {
		name string
		args args
		want *domain.Product
	}{
		{
			name: "should return domain Product record for the DB record",
			args: args{
				p: &Product{
					Id:                 productIDHex,
					Name:               "test name",
					SubscriptionPeriod: 1,
					Price:              10,
					TaxPercentage:      10,
				},
			},
			want: &domain.Product{
				ID:                 productIDHex.Hex(),
				Name:               "test name",
				SubscriptionPeriod: 1,
				Price:              10,
				TaxPercentage:      10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createDomainProductRecord(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDomainProductRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *MongoTestSuite) TestGetProduct() {
	mgoC := suite.TestContainer
	t := suite.T()
	client, err := connect(fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port))
	if err != nil {
		t.Fatal(err)
	}

	dbName := "testdb"
	productIDHex := primitive.NewObjectID()
	_, err = client.Database(dbName).Collection(productCollection).InsertOne(context.Background(), &Product{
		Id:                 productIDHex,
		Name:               "test name",
		SubscriptionPeriod: 1,
		Price:              10,
		TaxPercentage:      10,
	})

	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		client                     *mongo.Client
		dbName                     string
		ProductCollection          *mongo.Collection
		UserSubscriptionCollection *mongo.Collection
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Product
		wantErr bool
	}{
		{
			name: "should return record for given id",
			fields: fields{
				client:                     client,
				dbName:                     "testdb",
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				id:  productIDHex.Hex(),
			},
			want: []domain.Product{
				{
					ID:                 productIDHex.Hex(),
					Name:               "test name",
					SubscriptionPeriod: 1,
					Price:              10,
					TaxPercentage:      10,
				},
			},
			wantErr: false,
		},
		{
			name: "should return empty slice for id which is not present in the db",
			fields: fields{
				client:                     client,
				dbName:                     "testdb",
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				id:  primitive.NewObjectID().Hex(),
			},
			want:    []domain.Product{},
			wantErr: false,
		},
		{
			name: "should return error for invalid id",
			fields: fields{
				client:                     client,
				dbName:                     "testdb",
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				id:  "invalid id",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mongoDetails{
				client:                     tt.fields.client,
				dbName:                     tt.fields.dbName,
				ProductCollection:          tt.fields.ProductCollection,
				UserSubscriptionCollection: tt.fields.UserSubscriptionCollection,
			}
			got, err := m.GetProduct(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoDetails.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoDetails.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDomainProductRecordSl(t *testing.T) {
	productIDHex := primitive.NewObjectID()
	type args struct {
		p []Product
	}
	tests := []struct {
		name string
		args args
		want []domain.Product
	}{
		{
			name: "should return slice of domain products for valid db products",
			args: args{
				p: []Product{
					{
						Id:                 productIDHex,
						Name:               "test name",
						SubscriptionPeriod: 1,
						Price:              10,
						TaxPercentage:      10,
					},
				},
			},
			want: []domain.Product{
				{
					ID:                 productIDHex.Hex(),
					Name:               "test name",
					SubscriptionPeriod: 1,
					Price:              10,
					TaxPercentage:      10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createDomainProductRecordSl(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDomainProductRecordSl() = %v, want %v", got, tt.want)
			}
		})
	}
}
