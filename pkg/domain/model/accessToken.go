package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewAccessToken() *AccessToken {
	return &AccessToken{
		AccessToken: uuid.NewString(),
		TokenType:   "Bearer",
		ExpiresIn:   time.Now().Add(time.Minute).Unix(),
	}
}
