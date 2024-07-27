package main

import (
	"fmt"
	"os"

	authNServer "github.com/task4233/oauth/pkg/api/server/authentication"
	authZServer "github.com/task4233/oauth/pkg/api/server/authorization"
	"github.com/task4233/oauth/pkg/infra"
	authZUseCase "github.com/task4233/oauth/pkg/usecase/authorization"
	"golang.org/x/sync/errgroup"
)

const (
	AppServerPort            = 9000
	authorizationServerPort  = 9001
	authenticationServerPort = 9002
)

func main() {
	authStorage := infra.NewAuthorizationStorage()
	authZUC := authZUseCase.NewAuthUseCase(authStorage)
	authZSV := authZServer.NewAuthorization(authZUC)
	authNSV := authNServer.NewAuthentication()

	eg := &errgroup.Group{}

	eg.Go(func() error {
		return authZSV.Run(authorizationServerPort)
	})
	eg.Go(func() error {
		return authNSV.Run(authenticationServerPort)
	})

	if err := eg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
