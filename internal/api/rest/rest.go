package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ganeshdipdumbare/gymondo-subscription/internal/api"
	"github.com/ganeshdipdumbare/gymondo-subscription/internal/app"
)

const (
	nilArgErr   = "nil %v not allowed"
	emptyArgErr = "empty %v not allowed"
)

type apiDetails struct {
	app    app.App
	server *http.Server
}

func NewApi(a app.App, port string) (api.Api, error) {
	if a == nil {
		return nil, fmt.Errorf(nilArgErr, "app")
	}

	if port == "" {
		return nil, fmt.Errorf(emptyArgErr, "port")
	}

	api := &apiDetails{
		app: a,
	}

	router := api.setupRouter()
	api.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
		Handler: router,
	}

	return api, nil
}

// StartServer starts rest server and wait for kill signal to stop it gracefully
// otherwise returns error
func (a *apiDetails) StartServer() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// GracefulStopServer gracefully stops the rest server
func (a *apiDetails) GracefulStopServer() {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
