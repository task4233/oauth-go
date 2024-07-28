package model

type TokenType string

const (
	TokenTypeAccessToken  TokenType = "access_token"
	TokenTypeRefreshToken TokenType = "refresh_token"
)

// ref: https://datatracker.ietf.org/doc/html/rfc7662#section-2.2
type Introspect struct {
	Active    bool      // required
	Scope     string    // optional
	ClientID  string    // optional
	Username  string    // optional
	TokenType TokenType // optional
	Exp       int64     // optional, expiration time
	Iat       int64     // optional, issued at
	Nbf       int64     // optional, not to be used before
	Sub       string    // optional, subject of the token
	Aud       string    // optional, audience
	Jti       string    // optional, JWT ID
}
