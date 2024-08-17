package model

import (
	"slices"
)

// ref:
// - https://datatracker.ietf.org/doc/html/rfc6749#section-2.3
// - http://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#ClientAuthentication
type AuthMethod string

const (
	AuthMethodBasic AuthMethod = "client_secret_basic"
)

type Client interface {
	GetAuthMethod() AuthMethod
	GetID() string
	GetSecret() string
	GetLoginURL(string) string
	GetRedirectURIs() []string
	IsPublic() bool
	IsValidRedirectURI(string) bool
}

type ConfidentialClient struct {
	authMethod   AuthMethod
	id           string
	secretHash   string
	redirectURIs []string
}

func NewConfidentialClient(
	authMethod AuthMethod,
	id string,
	secret string,
	redirectURIs []string,
) *ConfidentialClient {
	return &ConfidentialClient{
		authMethod:   authMethod,
		id:           id,
		secretHash:   secret,
		redirectURIs: redirectURIs,
	}
}

func (c *ConfidentialClient) GetID() string {
	return c.id
}

func (c *ConfidentialClient) GetAuthMethod() AuthMethod {
	return c.authMethod
}

func (c *ConfidentialClient) GetSecret() string {
	return c.secretHash
}

func (c *ConfidentialClient) GetRedirectURIs() []string {
	return c.redirectURIs
}

func (c *ConfidentialClient) IsPublic() bool {
	return false
}

func (c *ConfidentialClient) IsValidRedirectURI(uri string) bool {
	return slices.Contains(c.redirectURIs, uri)
}

func (c *ConfidentialClient) GetLoginURL(authReqID string) string {
	return "login?id=" + authReqID
}
