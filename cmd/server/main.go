package main

import (
	"context"
	"sync"

	"github.com/task4233/oauth/api"
	"github.com/task4233/oauth/infra"
	"github.com/task4233/oauth/logger"
)

const (
	authenticationServerPort = 7070
	authorizationServerPort  = 8080
	resourceServerPort       = 9090
)

func main() {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	kvs := infra.NewKVS()
	clients := map[string]*api.Client{
		"test_client": {
			ClientID:     "test_client",
			ClientSecret: "test_client_secret",
			Scope:        "read write",
		},
	}

	authServer := api.NewAuthorizationServer(ctx, authorizationServerPort, clients, kvs)
	resourceServer := api.NewResourceServer(ctx, resourceServerPort, kvs)

	wg := &sync.WaitGroup{}

	// run authorization server
	wg.Add(1)
	go func() {
		log.Info("authorization server is running on port %d\n", authorizationServerPort)
		if err := authServer.Run(); err != nil {
			log.Error("failed to run authorization server: %v", err)
			return
		}
		wg.Done()
	}()

	// run resource server
	wg.Add(1)
	go func() {
		log.Info("resource server is running on port %d\n", resourceServerPort)
		if err := resourceServer.Run(); err != nil {
			log.Error("failed to run resource server: %v", err)
			return
		}
		wg.Done()
	}()

	wg.Wait()
}
