package model

import "time"

type AuthRequest struct {
	ID           string
	ClientID     string
	Code         string
	RedirectURI  string
	ResponseType string
	State        string
	Scope        string
	DisabledAt   time.Time
}

type TokenRequest struct {
	Code         string
	RedirectURI  string
	ClientID     string
	ClientSecret string
}
