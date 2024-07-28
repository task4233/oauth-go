package authorization

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/task4233/oauth/pkg/domain/model"
)

// ref: https://datatracker.ietf.org/doc/html/rfc7662#section-2.1
type IntrospectRequest struct {
	Token         string          `form:"token"`           // required
	TokenTypeHint model.TokenType `form:"token_type_hint"` // optional
}

func (r *IntrospectRequest) Validate() error {
	if r.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

// ref: https://datatracker.ietf.org/doc/html/rfc7662#section-2.2
type IntrospectResponse struct {
	Active    bool            `json:"active"`     // required
	Scope     string          `json:"scope"`      // optional
	ClientID  string          `json:"client_id"`  // optional
	Username  string          `json:"username"`   // optional
	TokenType model.TokenType `json:"token_type"` // optional
	Exp       int64           `json:"exp"`        // optional, expiration time
	Iat       int64           `json:"iat"`        // optional, issued at
	Nbf       int64           `json:"nbf"`        // optional, not to be used before
	Sub       string          `json:"sub"`        // optional, subject of the token
	Aud       string          `json:"aud"`        // optional, audience
	Jti       string          `json:"jti"`        // optional, JWT ID
}

func (s *Authorization) Introspect(w http.ResponseWriter, r *http.Request) {
	req := s.ParseIntrospectRequest(r)
	err := req.Validate()
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       InvalidRequest,
			Description: err.Error(),
		})
		return
	}

	res, err := s.authUC.Introspect(r.Context(), req.Token, req.TokenTypeHint)
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       ServerError,
			Description: err.Error(),
		})
		return
	}

	s.IntrospectResponse(w, r, res)
}

func (s *Authorization) IntrospectResponse(w http.ResponseWriter, r *http.Request, introspect *model.Introspect) {
	res := &IntrospectResponse{
		Active:    introspect.Active,
		Scope:     introspect.Scope,
		ClientID:  introspect.ClientID,
		Username:  introspect.Username,
		TokenType: introspect.TokenType,
		Exp:       introspect.Exp,
		Iat:       introspect.Iat,
		Nbf:       introspect.Nbf,
		Sub:       introspect.Sub,
		Aud:       introspect.Aud,
		Jti:       introspect.Jti,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       ServerError,
			Description: err.Error(),
		})
		return
	}
}

func (s *Authorization) ParseIntrospectRequest(r *http.Request) *IntrospectRequest {
	return &IntrospectRequest{
		Token:         r.FormValue("token"),
		TokenTypeHint: model.TokenType(r.FormValue("token_type_hint")),
	}
}
