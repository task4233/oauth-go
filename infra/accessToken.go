package infra

import (
	"context"
	"time"

	"github.com/task4233/oauth/domain"
	"github.com/task4233/oauth/infra/repository"
)

var _ repository.AccessToken = (*AccessToken)(nil)

const (
	accessTokenUserIDKey    = "accessToken_user_id"
	accessTokenClientIDKey  = "accessToken_client_id"
	accessTokenScopesKey    = "accessToken_scopes"
	accessTokenExpiresAtKey = "accessToken_expires_at"
)

type AccessToken struct {
	kvs repository.KVS
}

func NewAccessTokenRepository() *AccessToken {
	return &AccessToken{
		kvs: NewKVS(),
	}
}

func (c *AccessToken) Get(ctx context.Context, signature string) (*domain.AccessToken, error) {
	v, err := c.kvs.Get(signature)
	if err != nil {
		return nil, err
	}

	return &domain.AccessToken{
		Signature: signature,
		UserID:    v[accessTokenUserIDKey].(string),
		ClientID:  v[accessTokenClientIDKey].(string),
		Scopes:    v[accessTokenScopesKey].([]string),
		ExpiresAt: v[accessTokenExpiresAtKey].(time.Time),
	}, nil
}

func (c *AccessToken) Insert(ctx context.Context, accessToken *domain.AccessToken) error {
	if err := accessToken.Validate(); err != nil {
		return err
	}

	return c.kvs.Set(accessToken.Signature, map[string]any{
		"signature":             accessToken.Signature,
		accessTokenUserIDKey:    accessToken.UserID,
		accessTokenClientIDKey:  accessToken.ClientID,
		accessTokenScopesKey:    accessToken.Scopes,
		accessTokenExpiresAtKey: accessToken.ExpiresAt,
	})
}
