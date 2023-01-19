package mongodb

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/ganeshdipdumbare/gymondo-subscription/internal/db"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

type mongoTestContainer struct {
	Container testcontainers.Container
	Ip        string
	Port      string
}

type MongoTestSuite struct {
	suite.Suite
	TestContainer mongoTestContainer
}

func getMongoTestContainer(ctx context.Context) (*mongoTestContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
	}

	mgoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	ip, err := mgoC.Host(ctx)
	if err != nil {
		return nil, err
	}
	mongoPort, err := mgoC.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return nil, err

	}

	return &mongoTestContainer{
		Container: mgoC,
		Ip:        ip,
		Port:      mongoPort.Port(),
	}, nil
}

// SetupTest runs before every test
func (suite *MongoTestSuite) SetupTest() {
	testContainer, err := getMongoTestContainer(context.Background())
	if err != nil {
		log.Fatal("unable to get mongo test container")
	}
	suite.TestContainer = *testContainer
}

// TearDownTest runs after every test
func (suite *MongoTestSuite) TearDownTest() {
	suite.TestContainer.Container.Terminate(context.Background())
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}

func (suite *MongoTestSuite) TestNewMongoDB() {
	mgoC := suite.TestContainer
	t := suite.T()

	type args struct {
		uri    string
		dbName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error for empty uri",
			args: args{
				dbName: "gymondodb",
			},
			wantErr: true,
		},
		{
			name: "should return error for invalid uri",
			args: args{
				uri: fmt.Sprintf("mongodb://%s:%s", "invalidip", "invalidport"),
			},
			wantErr: true,
		},
		{
			name: "should return error for empty dbName",
			args: args{
				uri: fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port),
			},
			wantErr: true,
		},
		{
			name: "should return success for valid input args",
			args: args{
				uri:    fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port),
				dbName: "gymondodb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMongoDB(tt.args.uri, tt.args.dbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMongoDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (suite *MongoTestSuite) Test_connect() {
	mgoC := suite.TestContainer
	t := suite.T()

	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return success for valid input args",
			args: args{
				uri: fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port),
			},
			wantErr: false,
		},
		{
			name: "should return error for invalid input args",
			args: args{
				uri: fmt.Sprintf("mongodb://%s:%s", "invalidip", "invalidport"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := connect(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (suite *MongoTestSuite) TestDisconnect() {
	mgoC := suite.TestContainer
	t := suite.T()

	mongodb, err := NewMongoDB(fmt.Sprintf("mongodb://%s:%s", mgoC.Ip, mgoC.Port), "testdb")
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		database db.DB
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should return success for valid db",
			args: args{
				in0: context.Background(),
			},
			fields: fields{
				database: mongodb,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.database
			if err := m.Disconnect(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("mongoDetails.Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
