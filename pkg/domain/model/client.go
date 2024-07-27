package model

type Client interface {
	GetLoginURL(string) string
}

type ConfidentialClient struct {
}

func (c *ConfidentialClient) GetLoginURL(authReqID string) string {
	return "login?id=" + authReqID
}
