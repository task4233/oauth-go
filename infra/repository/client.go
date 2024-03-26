package repository

import (
	"context"

	"github.com/task4233/oauth/domain"
)

type Client interface {
	Get(ctx context.Context, id string) (*domain.Client, error)
}
