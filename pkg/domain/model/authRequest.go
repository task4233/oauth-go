package model

type AuthRequest struct {
	ID           string
	ClientID     string
	Code         string
	RedirectURI  string
	ResponseType string
	State        string
	Scope        string
}
