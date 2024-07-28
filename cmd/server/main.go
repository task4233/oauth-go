package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/task4233/oauth/pkg/api/client"
	authNServer "github.com/task4233/oauth/pkg/api/server/authentication"
	authZServer "github.com/task4233/oauth/pkg/api/server/authorization"
	resourceServer "github.com/task4233/oauth/pkg/api/server/resource"
	"github.com/task4233/oauth/pkg/infra"
	authZUseCase "github.com/task4233/oauth/pkg/usecase/authorization"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const (
	AppServerPort            = 9000
	authorizationServerPort  = 9001
	authenticationServerPort = 9002
	resourceServerPort       = 9003
)

func main() {
	oauthConfig := setupOAuthConfig()

	authStorage := infra.NewAuthorizationStorage()
	authZUC := authZUseCase.NewAuthUseCase(authStorage)
	authZSV := authZServer.NewAuthorization(authZUC)
	authNSV := authNServer.NewAuthentication()
	resourceSV := resourceServer.NewResource()
	appSV := client.NewApp(oauthConfig)

	eg := &errgroup.Group{}

	eg.Go(func() error {
		return appSV.Run(AppServerPort)
	})
	eg.Go(func() error {
		return authZSV.Run(authorizationServerPort)
	})
	eg.Go(func() error {
		return authNSV.Run(authenticationServerPort)
	})
	eg.Go(func() error {
		return resourceSV.Run(resourceServerPort)
	})

	if err := eg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func setupOAuthConfig() *oauth2.Config {
	authZServerBaseURL := "http://localhost:" + strconv.Itoa(authorizationServerPort)
	authZURL := authZServerBaseURL + "/authorize"
	tokenURL := authZServerBaseURL + "/token"

	clientBaseURL := "http://localhost:" + strconv.Itoa(AppServerPort)
	callbackURL := clientBaseURL + "/auth/callback"

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authZURL,
			TokenURL: tokenURL,
		},
		RedirectURL: callbackURL,
		Scopes:      []string{"openid", "profile", "email"},
	}

	return oauthConfig
}
