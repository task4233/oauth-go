package repository

import (
	"context"

	"github.com/task4233/oauth/domain"
)

type AccessToken interface {
	Get(ctx context.Context, signature string) (*domain.AccessToken, error)
	Insert(ctx context.Context, accessToken *domain.AccessToken) error
}
