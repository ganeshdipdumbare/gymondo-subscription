// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2022-06-30 11:27:48.521102 +0200 CEST m=+0.050108522

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "A REST server to manage user subscriptions of the products",
        "title": "Gymondo subscription api",
        "contact": {},
        "license": {},
        "version": "1.0"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/product": {
            "get": {
                "description": "return feteched  products",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product-api"
                ],
                "summary": "get all the products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.getAllProductsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    }
                }
            }
        },
        "/product/{id}": {
            "get": {
                "description": "return feteched  product record for input id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product-api"
                ],
                "summary": "get a product for given product id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.getProductByIdResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    }
                }
            }
        },
        "/subscription": {
            "post": {
                "description": "return created subscription record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription-api"
                ],
                "summary": "create a subscription for the user with given product",
                "parameters": [
                    {
                        "description": "create subscription request",
                        "name": "buySubscriptionRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.buySubscriptionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.buySubscriptionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    }
                }
            }
        },
        "/subscription/{id}": {
            "get": {
                "description": "return feteched  subscription record for input id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription-api"
                ],
                "summary": "get a subscription for given subscription id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.getSubscriptionByIDResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    }
                }
            }
        },
        "/subscription/{id}/changeStatus/{status}": {
            "put": {
                "description": "update subscription with given status and returns updated subscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription-api"
                ],
                "summary": "update subscription with given status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "active",
                            "cancel",
                            "pause"
                        ],
                        "type": "string",
                        "description": "status",
                        "name": "status",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.updateSubscriptionByIDResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/rest.errorRespose"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "rest.buySubscriptionRequest": {
            "type": "object",
            "required": [
                "email_id",
                "product_id"
            ],
            "properties": {
                "email_id": {
                    "type": "string"
                },
                "product_id": {
                    "type": "string"
                }
            }
        },
        "rest.buySubscriptionResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "product_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "tax": {
                    "type": "number"
                }
            }
        },
        "rest.errorRespose": {
            "type": "object",
            "properties": {
                "errorMessage": {
                    "type": "string"
                }
            }
        },
        "rest.getAllProductsResponse": {
            "type": "object",
            "properties": {
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/rest.getProductByIdResponse"
                    }
                }
            }
        },
        "rest.getProductByIdResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "subscription_period": {
                    "type": "integer"
                },
                "tax_percentage": {
                    "type": "number"
                }
            }
        },
        "rest.getSubscriptionByIDResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "pause_start_date": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "product_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "tax": {
                    "type": "number"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "rest.updateSubscriptionByIDResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "pause_start_date": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "product_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "tax": {
                    "type": "number"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
