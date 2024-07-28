package authorization

import (
	"fmt"
	"net/http"

	"github.com/task4233/oauth/pkg/domain/model"
	"github.com/task4233/oauth/pkg/repository"
	"github.com/task4233/oauth/pkg/usecase/authorization"
)

type Authorizer interface {
	Storage() repository.Storage
}

type Authorization struct {
	Authorizer
	authUC *authorization.AuthUseCase
}

func NewAuthorization(authUC *authorization.AuthUseCase) *Authorization {
	return &Authorization{
		authUC: authUC,
	}
}

func (s *Authorization) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/authorize", s.Authorize)
	mux.HandleFunc("/token", s.Token)
	mux.HandleFunc("/introspect", s.Introspect)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// ref: http://openid-foundation-japan.github.io/rfc6749.ja.html#code-authz-req
type AuthorizationRequest struct {
	Scope        string // required
	ResponseType string // required
	ClientID     string // required
	RedirectURI  string // required
	State        string // recommended
	Client       string // optional (not in RFC), this value is set after user login
}

func (r *AuthorizationRequest) Validate() error {
	if r.Scope == "" {
		return fmt.Errorf("scope is required")
	}
	if r.ResponseType != "code" {
		return fmt.Errorf("response_type must be code")
	}
	if r.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if r.RedirectURI == "" {
		return fmt.Errorf("redirect_uri is required")
	}
	return nil
}

func (r *AuthorizationRequest) ToModel() *model.AuthRequest {
	return &model.AuthRequest{
		ID:           r.Client,
		ClientID:     r.ClientID,
		RedirectURI:  r.RedirectURI,
		ResponseType: r.ResponseType,
		State:        r.State,
		Scope:        r.Scope,
	}
}

// ref: http://openid-foundation-japan.github.io/rfc6749.ja.html#code-authz-resp
type AuthorizationResponse struct {
	Code  string `form:"code"`  // required
	State string `form:"state"` // recommended
}

func (s *Authorization) Authorize(w http.ResponseWriter, r *http.Request) {
	req := s.ParseAuthorizeRequest(r)

	// before login
	if req.Client == "" {
		err := req.Validate()
		if err != nil {
			RequestError(w, r, &ErrorResponse{
				Error:       InvalidRequest,
				Description: err.Error(),
				ErrorURI:    req.RedirectURI,
				State:       req.State,
			})
			return
		}

		authReq, client, err := s.authUC.AuthorizeBeforeLogin(r.Context(), req.ToModel())
		if err != nil {
			RequestError(w, r, &ErrorResponse{
				Error:       ServerError,
				Description: err.Error(),
				ErrorURI:    req.RedirectURI,
				State:       req.State,
			})
			return
		}
		RedirectToLogin(w, r, client, authReq.ID)
		return
	}

	// after login
	authReq, client, err := s.authUC.AuthorizeAfterLogin(r.Context(), req.ToModel())
	if err != nil {
		RequestError(w, r, &ErrorResponse{
			Error:       ServerError,
			Description: err.Error(),
			ErrorURI:    req.RedirectURI,
			State:       req.State,
		})
		return
	}
	s.AuthResponseCode(w, r, authReq, client)
}

func (s *Authorization) AuthResponseCode(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, client model.Client) {
	res := &AuthorizationResponse{
		Code:  authReq.Code,
		State: authReq.State,
	}
	callback := fmt.Sprintf("%s?code=%s", authReq.RedirectURI, res.Code)
	if res.State != "" {
		callback += "&state=" + res.State
	}
	http.Redirect(w, r, callback, http.StatusFound)
}

func (s *Authorization) ParseAuthorizeRequest(r *http.Request) *AuthorizationRequest {
	res := &AuthorizationRequest{}
	res.Scope = r.FormValue("scope")
	res.ResponseType = r.FormValue("response_type")
	res.ClientID = r.FormValue("client_id")
	res.RedirectURI = r.FormValue("redirect_uri")
	res.State = r.FormValue("state")
	res.Client = r.URL.Query().Get("client")

	return res
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request, client model.Client, authReqID string) {
	loginURI := "http://localhost:9002/" + client.GetLoginURL(authReqID)
	http.Redirect(w, r, loginURI, http.StatusFound)
}
