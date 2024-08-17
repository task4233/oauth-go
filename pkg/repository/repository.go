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
	GetAccessToken(context.Context, string) (*model.AccessToken, error)

	GetClient(context.Context, string) (model.Client, error)
}

type Hasher interface {
	Compare(ctx context.Context, hash, data []byte) (bool, error)
	Hash(ctx context.Context, data []byte) ([]byte, error)
}
