package main

import (
	"context"
	"log"

	"github.com/task4233/oauth-go/api"
)

const port = 9002

func main() {
	ctx := context.Background()
	resource := api.NewResource(ctx, port)
	resource.Log.InfoContext(ctx, "start listening resource server on %d", port)
	if err := resource.Run(ctx); err != nil {
		log.Printf("failed resource.Run: %v", err.Error())
	}
}
