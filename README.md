[![Go](https://github.com/ganeshdipdumbare/gymondo-subscription/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/ganeshdipdumbare/gymondo-subscription/actions/workflows/go.yml)

<img align="right" width="180px" src="https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png">

# Gymondo Subscription API

The Gymondo Subscription API micro-service allows user to fetch products as well as subscribe to particular product.

## Contents
 - [Prerequisite](#prerequisite)
 - [Getting started](#getting-started)
 - [Description](#description)
    - [Use cases](#use-cases)
    - [API Operation](#api-operation)
    - [Technical details](#technical-details)
- [Improvements](#improvements)

## Prerequisite

1. go v1.18 (for running tests only)

2. docker

## Getting started
1. To run tests on the service
```sh
make test
```
2. To start the service
```sh
make start  
```
3. visit [swagger doc](http://localhost:8080/api/v1/swagger/index.html) in the browser and test APIs. Please refer `product` collection to get product ID or use API.

4. To stop the service
```sh  
make stop 
``` 
## Description
The microservice is used to fetch products and buy subscription with particular product.
## Use cases
1. User is able to fetch all the predefined products or a single product for given ID.
2. User is able to buy a single product results into starting the subscription for the period of `subscription period` of the product in terms of `month`.
3. User is able to pause the active subscription.User is able to activate the paused subscription again. The end date of subscription is extended for the time the subscription was paused.
4. User is able to cancel the active/paused subscription. User is not allowed to change the suscription status once the subscription is cancelled.

## API Operation
1. Fetch all the products 
```
[GET] /api/v1/product
```
2. Fetch product with given product ID 
```
[GET] /api/v1/product/:id
```
3. Buy a subscription for particular product
```
[POST] /api/v1/subscription
# sample body 
{
  "email_id": "test@test.com",
  "product_id": "62bac24b0bf33af1c877d97f"
}
```
4. Fetch subscription details for given subscription ID
```
[GET] /api/v1/subscription/:id
```
5. Change subscription status for given subscription ID
```
[PUT] /api/v1/subscription/:id/changeStatus/:status
```

## Technical details
- The service is written using clean code architecture which makes it modular and easy to maintain and test. These are the following layers  -
    - domain - Inner most layer, no external dependencies
    - app - consists of main business logic, dependent on domain only (must have 100% test coverage)
    - db - consists of db interface which provides db functions
        - URL - localhost:27017
        - DB - gymondodb
        - Product Collection - `product` created during migration at the start of service stores product records.
        - User Subscription Collection - `user_subscription` store user subscription records.
    - config - consists of functions crucial to start the service
    - migration - consists of files used in migration of `reference data`. In our case `product` data.  
    - api - the layer is used to communicate with the service. The new APIs like grpc or graphQL can be implemented in this layer by keeping other layers intact.
- The product data is migrated at the start of the service.

## Improvements
- Just a sample code, not as per system design which requires exact requirements
- subscription start and end date is `DateTime` to make it simpler for testing
- Subscription auto extend after end date 
- Decide which DB can be used as per the data and accordingly may need normalization.
- As of now, product name is stored in subscription details to make it simpler for testing.
- Use pagination for getting the products if the records in high quantity.
- Use of authentication/authorization for the user.
- Add more test cases
