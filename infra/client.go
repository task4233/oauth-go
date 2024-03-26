package infra

import (
	"context"

	"github.com/task4233/oauth/domain"
	"github.com/task4233/oauth/infra/repository"
)

var _ repository.Client = (*Client)(nil)

const (
	clientScopesKey     = "client_scopes"
	clientSecretHashKey = "client_secret_hash"
	clientNameKey       = "client_name"
)

type Client struct {
	kvs repository.KVS
}

func NewClientRepository() *Client {
	return &Client{
		kvs: NewKVS(),
	}
}

func (c *Client) Get(ctx context.Context, id string) (*domain.Client, error) {
	v, err := c.kvs.Get(id)
	if err != nil {
		return nil, err
	}

	return &domain.Client{
		ID:         id,
		SecretHash: v[clientSecretHashKey].(string),
		Name:       v[clientNameKey].(string),
		Scopes:     v[clientScopesKey].([]string),
	}, nil
}

func (c *Client) Insert(ctx context.Context, client *domain.Client) error {
	return c.kvs.Set(client.ID, map[string]any{
		"id":                client.ID,
		clientNameKey:       client.Name,
		clientScopesKey:     client.Scopes,
		clientSecretHashKey: client.SecretHash,
	})
}
