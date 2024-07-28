package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresIn    int64
	Scope        string // space-delimited, ref: https://datatracker.ietf.org/doc/html/rfc6749#section-3.3
}

func NewAccessToken(authReq *AuthRequest) *AccessToken {
	if authReq == nil {
		return nil
	}

	return &AccessToken{
		AccessToken: uuid.NewString(),
		TokenType:   "Bearer",
		ExpiresIn:   time.Now().Add(time.Minute).Unix(),
		Scope:       authReq.Scope,
	}
}
