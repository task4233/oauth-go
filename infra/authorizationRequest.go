package infra

import (
	"context"

	"github.com/task4233/oauth/domain"
	"github.com/task4233/oauth/infra/repository"
)

var _ repository.AuthorizationRequest = (*AuthorizationRequest)(nil)

const (
	authorizationRequestClientIDKey    = "auth_req_client_id"
	authorizationRequestStateKey       = "auth_req_state"
	authorizationRequestScopesKey      = "auth_req_scopes"
	authorizationRequestRedirectURIKey = "auth_req_redirect_uri"
)

type AuthorizationRequest struct {
	kvs repository.KVS
}

func NewAuthorizationRequestRepository() *AuthorizationRequest {
	return &AuthorizationRequest{
		kvs: NewKVS(),
	}
}

func (c *AuthorizationRequest) Get(ctx context.Context, id string) (*domain.AuthorizationRequest, error) {
	v, err := c.kvs.Get(id)
	if err != nil {
		return nil, err
	}

	return &domain.AuthorizationRequest{
		ID:          id,
		ClientID:    v[authorizationRequestClientIDKey].(string),
		State:       v[authorizationRequestStateKey].(string),
		Scopes:      v[authorizationRequestScopesKey].([]string),
		RedirectURI: v[authorizationRequestRedirectURIKey].(string),
	}, nil
}

func (c *AuthorizationRequest) Insert(ctx context.Context, authorizationRequest *domain.AuthorizationRequest) error {
	if err := authorizationRequest.Validate(); err != nil {
		return err
	}

	return c.kvs.Set(authorizationRequest.ID, map[string]any{
		"id":                               authorizationRequest.ID,
		authorizationRequestClientIDKey:    authorizationRequest.ClientID,
		authorizationRequestStateKey:       authorizationRequest.State,
		authorizationRequestScopesKey:      authorizationRequest.Scopes,
		authorizationRequestRedirectURIKey: authorizationRequest.RedirectURI,
	})
}
