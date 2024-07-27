package infra

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/task4233/oauth/pkg/domain/model"
)

var (
	ErrAuthReqNotFound = errors.New("authorization request not found")
	ErrClientNotFound  = errors.New("client not found")
)

type AuthorizationStorage struct {
	authReqKvs     map[string]*model.AuthRequest
	accessTokenKvs map[string]*model.AccessToken
	clientKvs      map[string]model.Client
}

func NewAuthorizationStorage() *AuthorizationStorage {
	return &AuthorizationStorage{
		authReqKvs:     make(map[string]*model.AuthRequest),
		accessTokenKvs: make(map[string]*model.AccessToken),
		clientKvs: map[string]model.Client{
			"dummy-client-id": &model.ConfidentialClient{
				ID: "dummy-client-id",
				RedirectURIs: []string{
					"http://localhost:9000/auth/callback",
				},
			},
		},
	}
}

func (s *AuthorizationStorage) GetAuthorizationRequest(ctx context.Context, id string) (*model.AuthRequest, error) {
	v, ok := s.authReqKvs[id]
	if !ok {
		return nil, ErrAuthReqNotFound
	}
	return v, nil
}

func (s *AuthorizationStorage) GetAuthorizationRequestByCode(ctx context.Context, code string) (*model.AuthRequest, error) {
	for _, v := range s.authReqKvs {
		if v.Code == code {
			return v, nil
		}
	}

	return nil, ErrAuthReqNotFound
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

func (s *AuthorizationStorage) DisableAuthorizationRequest(ctx context.Context, id string) error {
	_, ok := s.authReqKvs[id]
	if !ok {
		return ErrAuthReqNotFound
	}

	s.authReqKvs[id].DisabledAt = time.Now()
	return nil
}

func (s *AuthorizationStorage) CreateAccessToken(ctx context.Context, token *model.AccessToken) error {
	s.accessTokenKvs[token.AccessToken] = token
	return nil
}

func (s *AuthorizationStorage) GetClient(ctx context.Context, clientID string) (model.Client, error) {
	client, ok := s.clientKvs[clientID]
	if !ok {
		return nil, ErrClientNotFound
	}

	_, ok = client.(*model.ConfidentialClient)
	if !ok {
		return nil, ErrClientNotFound
	}

	return client, nil
}
