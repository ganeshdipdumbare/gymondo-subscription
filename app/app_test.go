package app

import (
	"context"
	"reflect"
	"testing"

	"github.com/ganeshdipdumbare/gymondo-subscription/db"
	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"github.com/ganeshdipdumbare/gymondo-subscription/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AppTestSuite struct {
	suite.Suite
	Database       *mocks.MockDB
	MockController *gomock.Controller
}

// SetupTest runs before every test
func (suite *AppTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.MockController = mockCtrl
	suite.Database = mocks.NewMockDB(mockCtrl)
}

// TearDownTest runs after every test
func (suite *AppTestSuite) TearDownTest() {
	suite.MockController.Finish()
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (suite *AppTestSuite) TestNewApp() {
	t := suite.T()

	type args struct {
		database db.DB
	}
	tests := []struct {
		name    string
		args    args
		want    App
		wantErr bool
	}{
		{
			name: "should return app when valid input db",
			args: args{
				database: suite.Database,
			},
			want: &appDetails{
				database: suite.Database,
			},
			wantErr: false,
		},
		{
			name: "should return error when nil input db",
			args: args{
				database: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApp(tt.args.database)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *AppTestSuite) Test_appDetails_GetProduct() {
	t := suite.T()

	database := suite.Database
	ctx := context.Background()
	id := "62bb4ecdba3bbe275f8c7788"
	productRecord := domain.Product{
		ID:                 id,
		Name:               "testname",
		SubscriptionPeriod: 1,
		Price:              10,
		TaxPercentage:      10,
	}

	gomock.InOrder(
		database.EXPECT().GetProduct(ctx, id).Return([]domain.Product{
			productRecord,
		}, nil).Times(1),
		database.EXPECT().GetProduct(ctx, "invalidid").Return(nil,
			db.InvalidArgErr).Times(1),
		database.EXPECT().GetProduct(ctx, "").Return([]domain.Product{
			productRecord,
		}, nil).Times(1),
	)
	type fields struct {
		database db.DB
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
			name: "should return product for valid input id",
			fields: fields{
				database: database,
			},
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: []domain.Product{
				productRecord,
			},
			wantErr: false,
		},
		{
			name: "should return error for invalid input id",
			fields: fields{
				database: database,
			},
			args: args{
				ctx: ctx,
				id:  "invalidid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return slice of products for empty product id",
			fields: fields{
				database: database,
			},
			args: args{
				ctx: ctx,
				id:  "",
			},
			want: []domain.Product{
				productRecord,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &appDetails{
				database: tt.fields.database,
			}
			got, err := a.GetProduct(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("appDetails.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appDetails.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *AppTestSuite) TestBuySubscription() {
	t := suite.T()

	database := suite.Database
	subscriptionId := "62bb4ecdba3bbe275f8c7788"
	productId := "62bb4ecdba3bbe275f8c7788"
	ctx := context.Background()
	subscriptionRecord := domain.UserSubscription{
		ID: subscriptionId,
	}
	productRecord := domain.Product{
		ID:                 productId,
		Name:               "testproduct",
		SubscriptionPeriod: 1,
		Price:              10,
		TaxPercentage:      10,
	}

	gomock.InOrder(
		database.EXPECT().GetProduct(gomock.Any(), gomock.AssignableToTypeOf(productRecord.ID)).Return([]domain.Product{
			productRecord,
		}, nil),
		database.EXPECT().SaveSubscription(gomock.Any(), gomock.AssignableToTypeOf(&subscriptionRecord)).Return(&subscriptionRecord, nil),
	)

	type fields struct {
		database db.DB
	}
	type args struct {
		ctx       context.Context
		productID string
		emailID   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.UserSubscription
		wantErr bool
	}{
		{
			name: "should return success for valid inputs",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				productID: productId,
				emailID:   "testmail@test.com",
			},
			wantErr: false,
		},
		{
			name: "should return error for empty inputs",
			fields: fields{
				database: database,
			},
			args: args{
				ctx:       ctx,
				productID: "",
				emailID:   "testmail@test.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &appDetails{
				database: tt.fields.database,
			}
			_, err := a.BuySubscription(tt.args.ctx, tt.args.productID, tt.args.emailID)
			if (err != nil) != tt.wantErr {
				t.Errorf("appDetails.BuySubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
