package main

import (
	"context"

	"github.com/task4233/oauth-go/api"
	"github.com/task4233/oauth-go/infra"
)

const (
	authorizationPort = 9001
	resourcePort      = 9002
)

func main() {
	ctx := context.Background()

	clients := []*api.Client{{
		ClientID:     "oauth-client-id-1",
		ClientSecret: "oauth-client-secret-1",
		RedirectURI:  []string{"http://localhost:9000/callback"},
		Scope:        "read write",
	}}
	resourceData := map[string]string{
		"name":        "protected resource",
		"description": "this data is protected by OAuth2.0",
	}

	kvs := infra.NewKVS()
	authorization := api.NewAuthorization(ctx, authorizationPort, clients, kvs)
	resource := api.NewResource(ctx, resourcePort, kvs, resourceData)

	// running authorization server
	go func() {
		authorization.Log.InfoContext(ctx, "start listening authorization server on %d", authorizationPort)
		if err := authorization.Run(ctx); err != nil {
			authorization.Log.WarnContext(ctx, "failed authorization.Run: %v", err.Error())
		}
	}()

	// running resource server
	go func() {
		resource.Log.InfoContext(ctx, "start listening resource server on %d", resourcePort)
		if err := resource.Run(ctx); err != nil {
			resource.Log.WarnContext(ctx, "failed resource.Run: %v", err.Error())
		}
	}()

}
