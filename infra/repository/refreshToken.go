package repository

import (
	"context"

	"github.com/task4233/oauth/domain"
)

type RefreshToken interface {
	Get(ctx context.Context, signature string) (*domain.RefreshToken, error)
	Insert(ctx context.Context, refreshToken *domain.RefreshToken) error
}
