package infra

import (
	"context"
	"time"

	"github.com/task4233/oauth/domain"
	"github.com/task4233/oauth/infra/repository"
)

var _ repository.RefreshToken = (*RefreshToken)(nil)

const (
	refreshTokenUserIDKey     = "refreshToken_user_id"
	refreshTokenClientIDKey   = "refreshToken_client_id"
	refreshTokenScopesKey     = "refreshToken_scopes"
	refreshTokenExpiresAtKey  = "refreshToken_expires_at"
	refreshTokenDisabledAtKey = "refreshToken_disabled_at"
)

type RefreshToken struct {
	kvs repository.KVS
}

func NewRefreshTokenRepository() *RefreshToken {
	return &RefreshToken{
		kvs: NewKVS(),
	}
}

func (c *RefreshToken) Get(ctx context.Context, signature string) (*domain.RefreshToken, error) {
	v, err := c.kvs.Get(signature)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshToken{
		Signature:  signature,
		UserID:     v[refreshTokenUserIDKey].(string),
		ClientID:   v[refreshTokenClientIDKey].(string),
		Scopes:     v[refreshTokenScopesKey].([]string),
		ExpiresAt:  v[refreshTokenExpiresAtKey].(time.Time),
		DisabledAt: v[refreshTokenDisabledAtKey].(time.Time),
	}, nil
}

func (c *RefreshToken) Insert(ctx context.Context, refreshToken *domain.RefreshToken) error {
	if err := refreshToken.Validate(); err != nil {
		return err
	}

	return c.kvs.Set(refreshToken.Signature, map[string]any{
		"signature":               refreshToken.Signature,
		refreshTokenUserIDKey:     refreshToken.UserID,
		refreshTokenClientIDKey:   refreshToken.ClientID,
		refreshTokenScopesKey:     refreshToken.Scopes,
		refreshTokenExpiresAtKey:  refreshToken.ExpiresAt,
		refreshTokenDisabledAtKey: refreshToken.DisabledAt,
	})
}
