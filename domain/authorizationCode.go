package domain

import "time"

type AuthorizationCode struct {
	Code        string
	UserID      string
	ClientID    string
	Scopes      []string
	RedirectURI string
	ExpiresAt   time.Time
	DisabledAt  time.Time
}
