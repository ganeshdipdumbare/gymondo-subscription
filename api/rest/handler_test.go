package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ganeshdipdumbare/gymondo-subscription/app"
	"github.com/ganeshdipdumbare/gymondo-subscription/domain"
	"github.com/ganeshdipdumbare/gymondo-subscription/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type HandlerTestSuite struct {
	suite.Suite
	App            *mocks.MockApp
	MockController *gomock.Controller
}

// SetupTest runs before every test
func (suite *HandlerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.MockController = mockCtrl
	suite.App = mocks.NewMockApp(mockCtrl)
}

// TearDownTest runs after every test
func (suite *HandlerTestSuite) TearDownTest() {
	suite.MockController.Finish()
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestGetProductByID() {
	t := suite.T()

	appInstance := suite.App
	productID := "62bc589278b49cee00f01421"
	productIDNotPresent := "62bc59417f1271f8e9c5e1c4"
	productRecord := domain.Product{
		ID:                 productID,
		Name:               "test name",
		SubscriptionPeriod: 1,
		Price:              10,
		TaxPercentage:      10,
	}

	gomock.InOrder(
		appInstance.EXPECT().GetProduct(gomock.Any(), productID).Return([]domain.Product{
			productRecord,
		}, nil).Times(1),
		appInstance.EXPECT().GetProduct(gomock.Any(), "invalidid").Return(nil,
			app.InvalidArgErr).Times(1),
		appInstance.EXPECT().GetProduct(gomock.Any(), productIDNotPresent).Return([]domain.Product{}, nil).Times(1),
	)

	api := &apiDetails{
		app: appInstance,
	}
	router := api.setupRouter()
	var errorResp errorRespose

	// success test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/"+productID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var v getProductByIdResponse
	json.NewDecoder(w.Body).Decode(&v)
	assert.DeepEqual(t, getProductByIdResponse{
		ID:                 productRecord.ID,
		Name:               productRecord.Name,
		SubscriptionPeriod: productRecord.SubscriptionPeriod,
		Price:              productRecord.Price,
		TaxPercentage:      productRecord.TaxPercentage,
	}, v)

	// invalid id in the input
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/product/"+"invalidid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	json.NewDecoder(w.Body).Decode(&errorResp)
	assert.DeepEqual(t, errorRespose{
		ErrorMessage: "invalid argument",
	}, errorResp)

	// valid id for which product not present
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/product/"+productIDNotPresent, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	json.NewDecoder(w.Body).Decode(&errorResp)
	assert.DeepEqual(t, errorRespose{
		ErrorMessage: "product not found for given id",
	}, errorResp)
}

func (suite *HandlerTestSuite) TestGetAllProducts() {
	t := suite.T()

	appInstance := suite.App
	productID := "62bc589278b49cee00f01421"
	productRecord := domain.Product{
		ID:                 productID,
		Name:               "test name",
		SubscriptionPeriod: 1,
		Price:              10,
		TaxPercentage:      10,
	}

	appInstance.EXPECT().GetProduct(gomock.Any(), "").Return([]domain.Product{
		productRecord,
	}, nil).Times(1)

	api := &apiDetails{
		app: appInstance,
	}
	router := api.setupRouter()

	// success test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/product", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var v getAllProductsResponse
	json.NewDecoder(w.Body).Decode(&v)
	assert.DeepEqual(t, getAllProductsResponse{
		Products: []getProductByIdResponse{{
			ID:                 productRecord.ID,
			Name:               productRecord.Name,
			SubscriptionPeriod: productRecord.SubscriptionPeriod,
			Price:              productRecord.Price,
			TaxPercentage:      productRecord.TaxPercentage,
		}}}, v)
}
