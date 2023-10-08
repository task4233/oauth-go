package main

import (
	"context"

	"github.com/task4233/oauth/logger"
)

const (
	clientServerPort        = 8000
	authorizationServerPort = 8080
	resourceServerPort      = 9090
)

func main() {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	clients := client{
		clientID:     "test_client",
		clientSecret: "test_client_secret",
		scope:        "read write",
	}

	clientServer := NewClientServer(ctx, clientServerPort, clients)

	log.Info("client server is running", "port", clientServerPort)
	if err := clientServer.Run(); err != nil {
		log.Error("failed to run client server", "error", err)
		return
	}
}
