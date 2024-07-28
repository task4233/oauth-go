package authorization

import (
	"context"
	"fmt"

	"github.com/task4233/oauth/pkg/domain/model"
	"github.com/task4233/oauth/pkg/repository"
)

type AuthUseCase struct {
	Storage repository.Storage
}

func NewAuthUseCase(storage repository.Storage) *AuthUseCase {
	return &AuthUseCase{
		Storage: storage,
	}
}

func (s *AuthUseCase) AuthorizeBeforeLogin(ctx context.Context, req *model.AuthRequest) (*model.AuthRequest, model.Client, error) {
	authReq, err := s.Storage.CreateAuthorizationRequest(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	client, err := s.Storage.GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, nil, err
	}

	return authReq, client, nil
}

func (s *AuthUseCase) AuthorizeAfterLogin(ctx context.Context, req *model.AuthRequest) (*model.AuthRequest, model.Client, error) {
	authReq, err := s.Storage.GetAuthorizationRequest(ctx, req.ID)
	if err != nil {
		return nil, nil, err
	}

	client, err := s.Storage.GetClient(ctx, authReq.ClientID)
	if err != nil {
		return nil, nil, err
	}

	authReq, err = s.Storage.GenerateAuthorizationCode(ctx, authReq)
	if err != nil {
		return nil, nil, err
	}

	return authReq, client, nil
}

func (s *AuthUseCase) Token(ctx context.Context, req *model.TokenRequest) (*model.AccessToken, error) {
	authReq, _, err := s.ValidateTokenRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	err = s.Storage.DisableAuthorizationRequest(ctx, authReq.ID)
	if err != nil {
		return nil, err
	}

	accessToken := model.NewAccessToken(authReq)
	err = s.Storage.CreateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s *AuthUseCase) ValidateTokenRequest(ctx context.Context, req *model.TokenRequest) (*model.AuthRequest, model.Client, error) {
	client, err := s.Storage.GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, nil, err
	}

	authReq, err := s.Storage.GetAuthorizationRequestByCode(ctx, req.Code)
	if err != nil {
		return nil, nil, fmt.Errorf("by code: %w", err)
	}

	if client.GetID() != authReq.ClientID {
		return nil, nil, fmt.Errorf("client_id is mismatched")
	}
	if !client.IsValidRedirectURI(authReq.RedirectURI) {
		return nil, nil, fmt.Errorf("redirect_uri is mismatched")
	}

	return authReq, client, nil
}

func (s *AuthUseCase) Introspect(ctx context.Context, token string, hint model.TokenType) (*model.Introspect, error) {
	switch hint {
	case model.TokenTypeAccessToken:
		return s.introspectAccessToken(ctx, token)
	default:
		return nil, fmt.Errorf("unsupported token type")
	}
}

func (s *AuthUseCase) introspectAccessToken(ctx context.Context, token string) (*model.Introspect, error) {
	accessToken, err := s.Storage.GetAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &model.Introspect{
		Active:    true,
		Scope:     accessToken.Scope,
		ClientID:  "", // TODO: consider how to get client_id
		TokenType: model.TokenTypeAccessToken,
		Exp:       accessToken.ExpiresIn,
	}, nil
}
