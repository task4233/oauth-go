package domain

import "time"

type RefreshToken struct {
	Signature  string
	UserID     string
	ClientID   string
	Scopes     []string
	ExpiresAt  time.Time
	DisabledAt time.Time
}

func (r *RefreshToken) Validate() error {
	return nil
}
