package authorization

import (
	"context"

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
