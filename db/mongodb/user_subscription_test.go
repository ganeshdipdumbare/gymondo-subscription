package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_createDBUserSubscriptionRecord(t *testing.T) {
	timeNow := time.Now()
	type args struct {
		us *domain.UserSubscription
	}
	tests := []struct {
		name    string
		args    args
		want    *UserSubscription
		wantErr bool
	}{
		{
			name: "should return error for nil input",
			args: args{
				us: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return record for valid input record",
			args: args{
				us: &domain.UserSubscription{
					CreatedAt:      timeNow,
					Email:          "test@gmail.com",
					ProductName:    "test product",
					Status:         domain.SubscriptionStatusActive,
					StartDate:      timeNow,
					UpdatedAt:      &timeNow,
					EndDate:        timeNow,
					Price:          10.0,
					Tax:            10.0,
					PauseStartDate: &timeNow,
				},
			},
			want: &UserSubscription{
				CreatedAt:      timeNow,
				Email:          "test@gmail.com",
				ProductName:    "test product",
				Status:         string(domain.SubscriptionStatusActive),
				StartDate:      timeNow,
				UpdatedAt:      &timeNow,
				EndDate:        timeNow,
				Price:          10.0,
				Tax:            10.0,
				PauseStartDate: &timeNow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createDBUserSubscriptionRecord(tt.args.us)
			if (err != nil) != tt.wantErr {
				t.Errorf("createDBUserSubscriptionRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDBUserSubscriptionRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDomainUserSubscriptionRecord(t *testing.T) {
	timeNow := time.Now()
	idHex := primitive.NewObjectID()

	type args struct {
		us *UserSubscription
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.UserSubscription
		wantErr bool
	}{
		{
			name: "should return error for nil input",
			args: args{
				us: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return record for valid input record",
			args: args{
				us: &UserSubscription{
					Id:             idHex,
					CreatedAt:      timeNow,
					Email:          "test@gmail.com",
					ProductName:    "test product",
					Status:         string(domain.SubscriptionStatusActive),
					StartDate:      timeNow,
					UpdatedAt:      &timeNow,
					EndDate:        timeNow,
					Price:          10.0,
					Tax:            10.0,
					PauseStartDate: &timeNow,
				},
			},
			want: &domain.UserSubscription{
				ID:             idHex.Hex(),
				CreatedAt:      timeNow,
				Email:          "test@gmail.com",
				ProductName:    "test product",
				Status:         domain.SubscriptionStatusActive,
				StartDate:      timeNow,
				UpdatedAt:      &timeNow,
				EndDate:        timeNow,
				Price:          10.0,
				Tax:            10.0,
				PauseStartDate: &timeNow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createDomainUserSubscriptionRecord(tt.args.us)
			if (err != nil) != tt.wantErr {
				t.Errorf("createDomainUserSubscriptionRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDomainUserSubscriptionRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *MongoTestSuite) TestSaveSubscription() {
	mgoC := suite.TestContainer
	t := suite.T()
	client, err := connect(fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port))
	if err != nil {
		t.Fatal(err)
	}
	dbName := "testdb"
	timeNow := time.Now()
	idHex := primitive.NewObjectID()

	us := &domain.UserSubscription{
		ID:             idHex.Hex(),
		CreatedAt:      timeNow,
		Email:          "test@gmail.com",
		ProductName:    "test product",
		Status:         domain.SubscriptionStatusActive,
		StartDate:      timeNow,
		UpdatedAt:      &timeNow,
		EndDate:        timeNow,
		Price:          10.0,
		Tax:            10.0,
		PauseStartDate: &timeNow,
	}

	type fields struct {
		client                     *mongo.Client
		dbName                     string
		ProductCollection          *mongo.Collection
		UserSubscriptionCollection *mongo.Collection
	}
	type args struct {
		ctx context.Context
		us  *domain.UserSubscription
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.UserSubscription
		wantErr bool
	}{
		{
			name: "should return error for nil input",
			fields: fields{
				client:                     client,
				dbName:                     dbName,
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				us:  nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return record for valid input",
			fields: fields{
				client:                     client,
				dbName:                     dbName,
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				us:  us,
			},
			want:    us,
			wantErr: false,
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
			got, err := m.SaveSubscription(tt.args.ctx, tt.args.us)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoDetails.SaveSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoDetails.SaveSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *MongoTestSuite) TestGetSubscriptionByID() {
	mgoC := suite.TestContainer
	t := suite.T()
	client, err := connect(fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port))
	if err != nil {
		t.Fatal(err)
	}
	dbName := "testdb"
	timeNow := time.Now().UTC()
	idHex := primitive.NewObjectID()

	us := &domain.UserSubscription{
		ID:             idHex.Hex(),
		CreatedAt:      timeNow,
		Email:          "test@gmail.com",
		ProductName:    "test product",
		Status:         domain.SubscriptionStatusActive,
		StartDate:      timeNow,
		UpdatedAt:      &timeNow,
		EndDate:        timeNow,
		Price:          10.0,
		Tax:            10.0,
		PauseStartDate: &timeNow,
	}

	m := &mongoDetails{
		client:                     client,
		dbName:                     dbName,
		ProductCollection:          client.Database(dbName).Collection(productCollection),
		UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
	}
	us, err = m.SaveSubscription(context.Background(), us)
	if err != nil {
		t.Fatal("unable to save subscription")
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
		want    *domain.UserSubscription
		wantErr bool
	}{
		{
			name: "should return error for invalid id",
			fields: fields{
				client:                     client,
				dbName:                     dbName,
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				id:  "invalidid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return record for valid id",
			fields: fields{
				client:                     client,
				dbName:                     dbName,
				ProductCollection:          client.Database(dbName).Collection(productCollection),
				UserSubscriptionCollection: client.Database(dbName).Collection(userSubscriptionCollection),
			},
			args: args{
				ctx: context.Background(),
				id:  us.ID,
			},
			want:    us,
			wantErr: false,
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
			_, err := m.GetSubscriptionByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoDetails.GetSubscriptionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
