package domain

import "time"

type AccessToken struct {
	Signature string
	UserID    string
	ClientID  string
	Scopes    []string
	ExpiresAt time.Time
}

func (c *AccessToken) Validate() error {
	return nil
}
