package repository

import (
	"context"

	"github.com/task4233/oauth/domain"
)

type AuthorizationCode interface {
	Get(ctx context.Context, code string) (*domain.AuthorizationCode, error)
	Insert(ctx context.Context, authCode *domain.AuthorizationCode) error
}
