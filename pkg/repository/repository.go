package repository

import (
	"context"

	"github.com/task4233/oauth/pkg/domain/model"
)

type Storage interface {
	AuthorizationStorage
}

type AuthorizationStorage interface {
	GetAuthorizationRequest(context.Context, string) (*model.AuthRequest, error)
	GetAuthorizationRequestByCode(context.Context, string) (*model.AuthRequest, error)
	GenerateAuthorizationCode(context.Context, *model.AuthRequest) (*model.AuthRequest, error)
	CreateAuthorizationRequest(context.Context, *model.AuthRequest) (*model.AuthRequest, error)
	DisableAuthorizationRequest(context.Context, string) error

	CreateAccessToken(context.Context, *model.AccessToken) error

	GetClient(context.Context, string) (model.Client, error)
}
