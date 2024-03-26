package domain

type Client struct {
	ID         string
	SecretHash string
	Name       string
	Scopes     []string
}
