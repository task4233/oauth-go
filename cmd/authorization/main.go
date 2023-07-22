package main

import (
	"context"
	"log"
	"os"

	"github.com/task4233/oauth-go/api"
)

const port = 9001

func main() {
	ctx := context.Background()

	client := []*api.Client{{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  []string{"http://localhost:9000/callback"},
		Scope:        "read write",
	}}

	auth := api.NewAuthorization(ctx, port, client)
	auth.Log.InfoContext(ctx, "start listening authorization server on %d", port)
	if err := auth.Run(ctx); err != nil {
		log.Printf("failed auth.Run: %v", err.Error())
	}
}
