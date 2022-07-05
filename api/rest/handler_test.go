package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	getProdApiPath := "/api/v1/product/"
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
	req, _ := http.NewRequest(http.MethodGet, getProdApiPath+productID, nil)
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
	req, _ = http.NewRequest(http.MethodGet, getProdApiPath+"invalidid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	json.NewDecoder(w.Body).Decode(&errorResp)
	assert.DeepEqual(t, errorRespose{
		ErrorMessage: "invalid argument",
	}, errorResp)

	// valid id for which product not present
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, getProdApiPath+productIDNotPresent, nil)
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

	gomock.InOrder(
		appInstance.EXPECT().GetProduct(gomock.Any(), "").Return([]domain.Product{
			productRecord,
		}, nil).Times(1),
		appInstance.EXPECT().GetProduct(gomock.Any(), "").Return(nil, app.NotFoundErr).Times(1),
	)

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

	// error while getting product
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/product", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func (suite *HandlerTestSuite) TestBuySubscription() {
	t := suite.T()

	appInstance := suite.App
	productID := "62bc589278b49cee00f01421"
	subscriptionID := "62bc589278b49cee00f01421"
	emailID := "test@test.com"

	gomock.InOrder(
		appInstance.EXPECT().BuySubscription(gomock.Any(), productID, emailID).Return(&domain.UserSubscription{
			ID: subscriptionID,
		}, nil).Times(1),

		appInstance.EXPECT().BuySubscription(gomock.Any(), "invalidid", emailID).Return(nil, app.InvalidArgErr).Times(1),
	)

	api := &apiDetails{
		app: appInstance,
	}
	router := api.setupRouter()

	// success test
	w := httptest.NewRecorder()
	body := strings.NewReader(`{
		"product_id":"62bc589278b49cee00f01421",
		"email_id":"test@test.com"
	}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/subscription", body)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// invalid input email
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
		"product_id":"62bc589278b49cee00f01421",
		"email_id":"test"
	}`)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/subscription", body)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// invalid json body
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
		"product_id":"62bc589278b49cee00f01421"
		"email_id":"test"
	}`)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/subscription", body)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// buy subscription record returns error
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
		"product_id":"invalidid",
		"email_id":"test@test.com"
	}`)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/subscription", body)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func (suite *HandlerTestSuite) TestGetSubscriptionByID() {
	t := suite.T()

	appInstance := suite.App
	subscriptionID := "62bc589278b49cee00f01421"
	notFoundSubscriptionID := "62bc589278b49cee00f01421"

	gomock.InOrder(
		appInstance.EXPECT().GetSubscriptionByID(gomock.Any(), subscriptionID).Return(&domain.UserSubscription{
			ID: subscriptionID,
		}, nil).Times(1),

		appInstance.EXPECT().GetSubscriptionByID(gomock.Any(), "invalidid").Return(nil, app.InvalidArgErr).Times(1),

		appInstance.EXPECT().GetSubscriptionByID(gomock.Any(), notFoundSubscriptionID).Return(nil, app.NotFoundErr).Times(1),
	)

	api := &apiDetails{
		app: appInstance,
	}
	router := api.setupRouter()

	// success test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/subscription/"+subscriptionID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// invalid id test
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/subscription/"+"invalidid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// subscription not found for id test
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/subscription/"+notFoundSubscriptionID, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func (suite *HandlerTestSuite) TestUpdateSubscriptionStatusByID() {
	t := suite.T()

	appInstance := suite.App
	subscriptionID := "62bc589278b49cee00f01421"
	notFoundSubscriptionID := "62bc589278b49cee00f01421"

	gomock.InOrder(
		appInstance.EXPECT().UpdateSubscriptionStatusByID(gomock.Any(), subscriptionID, domain.SubscriptionStatusPaused).Return(&domain.UserSubscription{
			ID: subscriptionID,
		}, nil).Times(1),

		appInstance.EXPECT().UpdateSubscriptionStatusByID(gomock.Any(), notFoundSubscriptionID, domain.SubscriptionStatusPaused).Return(nil, app.NotFoundErr).Times(1),

		appInstance.EXPECT().UpdateSubscriptionStatusByID(gomock.Any(), "invalidID", domain.SubscriptionStatusPaused).Return(nil, app.InvalidArgErr).Times(1),
	)

	api := &apiDetails{
		app: appInstance,
	}
	router := api.setupRouter()

	// success test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/subscription/"+subscriptionID+"/changeStatus/pause", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// call with invalid status
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/subscription/"+subscriptionID+"/changeStatus/invalidstatus", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// subscription not found test
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/subscription/"+notFoundSubscriptionID+"/changeStatus/pause", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// call with invalid id test
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/api/v1/subscription/invalidID"+"/changeStatus/pause", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
