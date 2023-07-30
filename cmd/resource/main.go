package main

import (
	"context"
	"log"

	"github.com/task4233/oauth-go/api"
	"github.com/task4233/oauth-go/infra"
)

const port = 9002

func main() {
	ctx := context.Background()

	kvs := infra.NewKVS()
	resourceData := map[string]string{
		"name":        "protected resource",
		"description": "this data is protected by OAuth2.0",
	}
	resource := api.NewResource(ctx, port, kvs, resourceData)
	resource.Log.InfoContext(ctx, "start listening resource server on %d", port)
	if err := resource.Run(ctx); err != nil {
		log.Printf("failed resource.Run: %v", err.Error())
	}
}
