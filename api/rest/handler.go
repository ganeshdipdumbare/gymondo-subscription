package rest

import (
	"errors"
	"net/http"

	"github.com/ganeshdipdumbare/gymondo-subscription/app"
	docs "github.com/ganeshdipdumbare/gymondo-subscription/docs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	validate *validator.Validate
)

type getProductByIdResponse struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	SubscriptionPeriod uint    `json:"subscription_period"`
	Price              float64 `json:"price"`
	TaxPercentage      float64 `json:"tax_percentage"`
}

type getAllProductsResponse struct {
	Products []getProductByIdResponse `json:"products"`
}

type errorRespose struct {
	ErrorMessage string `json:"errorMessage"`
}

func createErrorResponse(c *gin.Context, code int, message string) {
	c.IndentedJSON(code, &errorRespose{
		ErrorMessage: message,
	})
}

func (api *apiDetails) setupRouter() *gin.Engine {
	validate = validator.New()

	apiV1 := "/api/v1"
	docs.SwaggerInfo.BasePath = apiV1

	r := gin.Default()
	v1group := r.Group(apiV1)
	v1group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1group.GET("/product/:id", api.getProductByID)
	v1group.GET("/product", api.getAllProducts)

	return r
}

// @BasePath /api/v1

// getProductByID godoc
// @Summary get a product for given product id
// @Description return feteched  product record for input id
// @Tags product-api
// @Accept  json
// @Produce  json
// @Param id path string true "product ID"
// @Success 200 {object} rest.getProductByIdResponse
// @Failure 404 {object} rest.errorRespose
// @Failure 400 {object} rest.errorRespose
// @Failure 500 {object} rest.errorRespose
// @Router /product/{id} [get]
func (api *apiDetails) getProductByID(c *gin.Context) {
	productID := c.Params.ByName("id")
	if productID == "" {
		createErrorResponse(c, http.StatusBadRequest, "param id cannot be empty")
		return
	}

	products, err := api.app.GetProduct(c, productID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, app.InvalidArgErr) {
			statusCode = http.StatusBadRequest
		}
		createErrorResponse(c, statusCode, err.Error())
		return
	}

	if len(products) == 0 {
		createErrorResponse(c, http.StatusNotFound, "product not found for given id")
		return
	}

	product := products[0]
	c.IndentedJSON(http.StatusOK, &getProductByIdResponse{
		ID:                 product.ID,
		Name:               product.Name,
		SubscriptionPeriod: product.SubscriptionPeriod,
		Price:              product.Price,
		TaxPercentage:      product.TaxPercentage,
	})
	c.Done()
}

// getAllProducts godoc
// @Summary get all the products
// @Description return feteched  products
// @Tags product-api
// @Accept  json
// @Produce  json
// @Success 200 {object} rest.getAllProductsResponse
// @Failure 500 {object} rest.errorRespose
// @Router /product [get]
func (api *apiDetails) getAllProducts(c *gin.Context) {
	products, err := api.app.GetProduct(c, "")
	if err != nil {
		createErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	respProducts := getAllProductsResponse{}
	for _, v := range products {
		respProducts.Products = append(respProducts.Products, getProductByIdResponse{
			ID:                 v.ID,
			Name:               v.Name,
			SubscriptionPeriod: v.SubscriptionPeriod,
			Price:              v.Price,
			TaxPercentage:      v.TaxPercentage,
		})
	}
	c.IndentedJSON(http.StatusOK, &respProducts)
	c.Done()
}
