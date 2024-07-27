package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/task4233/oauth/pkg/domain/model"
)

var ErrNotFound = errors.New("not found")

type AuthorizationStorage struct {
	authReqKvs map[string]*model.AuthRequest
}

func NewAuthorizationStorage() *AuthorizationStorage {
	return &AuthorizationStorage{
		authReqKvs: make(map[string]*model.AuthRequest),
	}
}

func (s *AuthorizationStorage) GetAuthorizationRequest(ctx context.Context, id string) (*model.AuthRequest, error) {
	v, ok := s.authReqKvs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (s *AuthorizationStorage) CreateAuthorizationRequest(ctx context.Context, req *model.AuthRequest) (*model.AuthRequest, error) {
	req.ID = uuid.NewString()
	s.authReqKvs[req.ID] = req
	return req, nil
}

func (s *AuthorizationStorage) GenerateAuthorizationCode(ctx context.Context, req *model.AuthRequest) (*model.AuthRequest, error) {
	req.Code = uuid.NewString()
	s.authReqKvs[req.ID] = req
	return req, nil
}

func (s *AuthorizationStorage) GetClient(ctx context.Context, clientID string) (model.Client, error) {
	confClient := &model.ConfidentialClient{}

	return confClient, nil
}
