package repository

import (
	"context"

	"github.com/task4233/oauth/domain"
)

type AuthorizationRequest interface {
	Get(ctx context.Context, id string) (*domain.AuthorizationRequest, error)
	Insert(ctx context.Context, session *domain.AuthorizationRequest) error
}
