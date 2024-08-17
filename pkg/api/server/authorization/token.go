package authorization

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/task4233/oauth/pkg/domain/model"
)

// ref: http://openid-foundation-japan.github.io/rfc6749.ja.html#token-req
type AccessTokenRequest struct {
	GrantType    string // required
	Code         string // required
	RedirectURI  string // required
	ClientID     string // required
	ClientSecret string // optional
}

func (r *AccessTokenRequest) Validate() error {
	if r.GrantType != "authorization_code" {
		return fmt.Errorf("grant_type must be authorization_code")
	}
	if r.Code == "" {
		return fmt.Errorf("code is required")
	}
	if r.RedirectURI == "" {
		return fmt.Errorf("redirect_uri is required")
	}
	if r.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}

	return nil
}

func (r *AccessTokenRequest) ToModel() *model.TokenRequest {
	return &model.TokenRequest{
		Code:         r.Code,
		RedirectURI:  r.RedirectURI,
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
	}
}

// ref: http://openid-foundation-japan.github.io/rfc6749.ja.html#token-response
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`  // required
	TokenType    string `json:"token_type"`    // required
	ExpiresIn    int64  `json:"expires_in"`    // recommended
	RefreshToken string `json:"refresh_token"` // optional
	Scope        string `json:"scope"`         // optional
}

func (s *Authorization) Token(w http.ResponseWriter, r *http.Request) {
	req := s.ParseAccessTokenRequest(r)
	err := req.Validate()
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       InvalidRequest,
			Description: err.Error(),
			ErrorURI:    req.RedirectURI,
		})
		return
	}

	accessToken, err := s.authUC.Token(r.Context(), req.ToModel())
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       ServerError,
			Description: err.Error(),
			ErrorURI:    req.RedirectURI,
		})
		return
	}

	s.TokenResponse(w, r, req, accessToken)
}

func (s *Authorization) TokenResponse(w http.ResponseWriter, r *http.Request, req *AccessTokenRequest, accessToken *model.AccessToken) {
	res := &AccessTokenResponse{
		AccessToken:  accessToken.AccessToken,
		TokenType:    accessToken.TokenType,
		ExpiresIn:    accessToken.ExpiresIn,
		RefreshToken: accessToken.RefreshToken,
		Scope:        accessToken.Scope,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       ServerError,
			Description: err.Error(),
			ErrorURI:    req.RedirectURI,
		})
		return
	}
}

func (s *Authorization) ParseAccessTokenRequest(r *http.Request) *AccessTokenRequest {
	req := &AccessTokenRequest{}
	req.GrantType = r.FormValue("grant_type")
	req.Code = r.FormValue("code")
	req.RedirectURI = r.FormValue("redirect_uri")

	clientID, clientSecret, ok := r.BasicAuth()
	if ok {
		req.ClientID = clientID
		req.ClientSecret = clientSecret
	} else {
		// ref: https://datatracker.ietf.org/doc/html/rfc6749#section-2.3.1
		req.ClientID = r.FormValue("client_id")
		req.ClientSecret = r.FormValue("client_secret")
	}

	return req
}
