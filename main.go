package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/api/rest"
	"github.com/ganeshdipdumbare/gymondo-subscription/app"
	"github.com/ganeshdipdumbare/gymondo-subscription/config"
	"github.com/ganeshdipdumbare/gymondo-subscription/db/mongodb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Gymondo Subscription API
// @version 1.0
// @description A REST server to manage user subscriptions of the products
// NewApi creates new api instance, otherwise returns error
func main() {
	// migrate reference data - product collection
	m, err := migrate.New(
		config.Get().MigrationFilesPath,
		config.Get().MongoUri+"/"+config.Get().MongoDb)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	// complete migration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	database, err := mongodb.NewMongoDB(config.Get().MongoUri, config.Get().MongoDb)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Disconnect(ctx)

	subscriptionApp, err := app.NewApp(database)
	if err != nil {
		log.Fatal(err)
	}

	restApi, err := rest.NewApi(subscriptionApp, config.Get().Port)
	if err != nil {
		log.Fatal(err)
	}
	restApi.StartServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	restApi.GracefulStopServer()
}
