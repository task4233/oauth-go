package domain

type AuthorizationRequest struct {
	ID          string
	ClientID    string
	State       string
	Scopes      []string
	RedirectURI string
}

func (s *AuthorizationRequest) Validate() error {
	if s == nil {
		return ErrIsNil
	}

	if s.ID == "" {
		return NewErrEmpty("id")
	}
	if s.ClientID == "" {
		return NewErrEmpty("client_id")
	}
	if s.State == "" {
		return NewErrEmpty("state")
	}
	if len(s.Scopes) == 0 {
		return NewErrEmpty("scopes")
	}

	return nil
}
