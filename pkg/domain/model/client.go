package model

import (
	"slices"
)

type Client interface {
	GetID() string
	GetLoginURL(string) string
	GetRedirectURIs() []string
	IsValidRedirectURI(string) bool
}

type ConfidentialClient struct {
	ID           string
	RedirectURIs []string
}

func (c *ConfidentialClient) GetID() string {
	return c.ID
}

func (c *ConfidentialClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *ConfidentialClient) IsValidRedirectURI(uri string) bool {
	return slices.Contains(c.RedirectURIs, uri)
}

func (c *ConfidentialClient) GetLoginURL(authReqID string) string {
	return "login?id=" + authReqID
}
