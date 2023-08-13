package main

import (
	"context"

	"github.com/task4233/oauth/api"
	"github.com/task4233/oauth/logger"
)

const (
	authorizationServerPort = 8080
	resourceServerPort      = 9090
)

func main() {
	ctx := context.Background()
	log := logger.FromContext(ctx)

	authServer := api.NewAuthorizationServer(authorizationServerPort)
	resourceServer := api.NewResourceServer(resourceServerPort)

	// run authorization server
	go func() {
		if err := authServer.Run(); err != nil {
			log.Error("failed to run authorization server: %v", err)
			return
		}
	}()

	// run resource server
	go func() {
		if err := resourceServer.Run(); err != nil {
			log.Error("failed to run resource server: %v", err)
			return
		}
	}()
}
