package main

import (
	"context"
	"sync"

	"github.com/task4233/oauth/api"
	"github.com/task4233/oauth/infra"
	"github.com/task4233/oauth/logger"
)

const (
	authorizationServerPort = 8080
	resourceServerPort      = 9090
)

func main() {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	kvs := infra.NewKVS()

	authServer := api.NewAuthorizationServer(authorizationServerPort, kvs, log)
	resourceServer := api.NewResourceServer(resourceServerPort, kvs, log)

	wg := &sync.WaitGroup{}

	// run authorization server
	wg.Add(1)
	go func() {
		log.Info("authorization server is running...", "port", authorizationServerPort)
		if err := authServer.Run(); err != nil {
			log.Error("failed to run authorization server", "error", err)
			return
		}
		wg.Done()
	}()

	// run resource server
	wg.Add(1)
	go func() {
		log.Info("resource server is running...", "port", resourceServerPort)
		if err := resourceServer.Run(); err != nil {
			log.Error("failed to run resource server", "error", err)
			return
		}
		wg.Done()
	}()

	wg.Wait()
}
