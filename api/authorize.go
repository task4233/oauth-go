package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/task4233/oauth-go/common"
	"github.com/task4233/oauth-go/infra/repository"
	"github.com/task4233/oauth-go/logger"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

// Client is a client application which is defined in RFC6749.
type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  []string
	Scope        string
}

// Code is a code which is defined in RFC6749.
type Code struct {
	AuthorizationEndpointRequest url.Values
	Scopes                       []string
	UserID                       string
}

// Authorization provides features for an authorization server of OAuth 2.0 which is defined in RFC6749.
// ref: https://datatracker.ietf.org/doc/html/rfc6749
type Authorization struct {
	AuthorizationEndpoint string
	kvs                   repository.KVS
	clients               map[string]*Client
	codes                 map[string]*Code
	requests              map[string]url.Values
	srv                   http.Server
	TokenEndpoint         string
	Log                   *slog.Logger
}

func NewAuthorization(
	ctx context.Context,
	port int,
	clients []*Client,
	kvs repository.KVS,
) Authorization {
	a := Authorization{
		Log: logger.FromContext(ctx),
	}
	a.kvs = kvs
	a.clients = map[string]*Client{}
	for _, c := range clients {
		a.clients[c.ClientID] = c
	}
	a.srv.Addr = fmt.Sprintf(":%d", port)
	a.srv.Handler = a.route()
	a.requests = map[string]url.Values{}
	a.codes = map[string]*Code{}

	return a
}

func (s *Authorization) Run(ctx context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *Authorization) route() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/approve", s.approve)
	r.HandleFunc("/authorize", s.authorize)
	r.HandleFunc("/token", s.token)

	return r
}

type AuthorizeResponse struct {
	Client Client `json:"client"`
	State  string `json:"state"`
	ReqID  string `json:"req_id"`
	Scope  string `json:"scope"`
}

// authorize provides "Authorization Endpoint".
// ref: https://datatracker.ietf.org/doc/html/rfc6749#section-3.1
func (s *Authorization) authorize(w http.ResponseWriter, r *http.Request) {
	s.Log.InfoContext(r.Context(), "GET /authorize is called")

	// get client with a clientID contained in query parameter in request
	clientID := r.URL.Query().Get("client_id")
	client, ok := s.clients[clientID]
	if !ok {
		msg := fmt.Sprintf("client_id is invalid: %s", clientID)
		s.Log.WarnContext(r.Context(), msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, msg)
		return
	}

	// check redirect_uri
	redirectURI := r.URL.Query().Get("redirect_uri")
	if !slices.Contains(client.RedirectURI, redirectURI) {
		msg := fmt.Sprintf("redirect_uri is invalid: %s, expected: %#v", redirectURI, client.RedirectURI)
		s.Log.WarnContext(r.Context(), msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, msg)
		return
	}

	// check scope
	reqScope := strings.Split(r.URL.Query().Get("scope"), " ")
	clientScope := strings.Split(client.Scope, " ")
	if !common.AreTwoUnorderedSlicesSame(reqScope, clientScope) {
		redirectURI, err := constructURIWithQueries(redirectURI, map[string]string{"error": "invalid_scope"})
		if err != nil {
			s.Log.WarnContext(r.Context(), err.Error())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		http.Redirect(w, r, redirectURI, http.StatusFound)
		return
	}

	// if all checks are passed, redirect to redirect_uri with code and state.
	reqID := uuid.New().String()
	s.requests[reqID] = r.URL.Query()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AuthorizeResponse{
		Client: *client,
		State:  r.Form.Get("state"),
		ReqID:  reqID,
		Scope:  client.Scope,
	})

}

func (s *Authorization) approve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.Log.InfoContext(ctx, "GET /approve is called")

	reqID := r.URL.Query().Get("req_id")
	req, ok := s.requests[reqID]
	if !ok {
		msg := fmt.Sprintf("no matched request with req_id: %s", reqID)
		s.Log.WarnContext(ctx, msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, msg)
		return
	}

	redirectURI := r.URL.Query().Get("redirect_uri")
	// if r.FormValue("approve") != "true" {
	// 	redirectURI, err := constructURIWithQueries(redirectURI, map[string]string{"error": "access_denied"})
	// 	if err != nil {
	// 		s.Log.WarnContext(ctx, err.Error())
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		fmt.Fprint(w, err.Error())
	// 		return
	// 	}
	// 	http.Redirect(w, r, redirectURI, http.StatusFound)
	// 	return
	// }

	switch req.Get("response_type") {
	case "code":
		code := uuid.New().String()
		userID := r.FormValue("user")

		// TODO: maybe fix this: remove prefix scope_ from scope?
		scope := strings.Split(r.FormValue("scope"), " ")
		client := s.clients[req.Get("client_id")]
		cScope := strings.Split(client.Scope, " ")
		if !common.AreTwoUnorderedSlicesSame(cScope, scope) {
			redirectURI, err := constructURIWithQueries(redirectURI, map[string]string{"error": "invalid_scope"})
			if err != nil {
				s.Log.WarnContext(ctx, err.Error())
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err.Error())
				return
			}
			http.Redirect(w, r, redirectURI, http.StatusFound)
			return
		}

		s.codes[code] = &Code{
			AuthorizationEndpointRequest: req,
			Scopes:                       scope,
			UserID:                       userID,
		}

		redirectURI, err := constructURIWithQueries(redirectURI, map[string]string{
			"code":  code,
			"state": req.Get("state"),
		})
		if err != nil {
			redirectURI, err = constructURIWithQueries(redirectURI, map[string]string{"error": "invalid_scope"})
			s.Log.WarnContext(ctx, err.Error())
			http.Redirect(w, r, redirectURI, http.StatusFound)
			return
		}
		http.Redirect(w, r, redirectURI, http.StatusFound)
	default:
		redirectURI, err := constructURIWithQueries(redirectURI, map[string]string{"error": "unsupported_response_type"})
		if err != nil {
			s.Log.WarnContext(ctx, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		http.Redirect(w, r, redirectURI, http.StatusFound)
	}
}

func (s *Authorization) token(w http.ResponseWriter, r *http.Request) {
	var clientID, clientSecret string
	var err error

	auth := r.Header.Get("Authorization")
	if auth != "" {
		// check the auth header
		clientID, clientSecret, err = parseBasicAuth(auth)
		if err != nil {
			s.Log.WarnContext(r.Context(), err.Error())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
	} else {
		if clientID != "" {
			s.Log.WarnContext(r.Context(), "client attempted to authenticate with multiple methods")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_client"})
			return
		}

		// check the body
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}

	client, ok := s.clients[clientID]
	if !ok {
		s.Log.WarnContext(r.Context(), fmt.Sprintf("client_id: %s is invalid", clientID))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_client"})
		return
	}
	if client.ClientSecret != clientSecret {
		s.Log.WarnContext(r.Context(), fmt.Sprintf("client_secret is invalid, expected %s got %s", client.ClientSecret, clientSecret))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_client"})
		return
	}

	switch r.FormValue("grant_type") {
	case "authorization_code":
		code := s.codes[r.FormValue("code")]
		if code == nil {
			s.Log.WarnContext(r.Context(), fmt.Sprintf("code: %s is invalid", r.FormValue("code")))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_grant"})
			return
		}
		expectedClientID := code.AuthorizationEndpointRequest.Get("client_id")
		if expectedClientID != clientID {
			s.Log.WarnContext(r.Context(), fmt.Sprintf("client_id is mismatch, expected %s got %s", expectedClientID, clientID))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_grant"})
			return
		}

		// TODO: replace the way to make accessToken because uuid is not suitable for accessToken
		accessToken := uuid.New().String()
		cScope := strings.Join(code.Scopes, " ")

		// TODO: insert accessToken, clientID, clientScope into kvs\
		vv := map[string]string{
			"access_token": accessToken,
			"client_id":    clientID,
			"scope":        cScope,
		}
		s.kvs.Set(accessToken, vv)

		tokenResponse := map[string]string{
			"access_token": accessToken,
			"token_type":   "Bearer",
			"scope":        cScope,
		}
		s.Log.InfoContext(r.Context(), fmt.Sprintf("access token is issued: %#v", tokenResponse))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokenResponse)
	case "refresh_token":
		panic("implement refresh token grant type")
	default:
		msg := fmt.Sprintf("unknown grant_type: %s", r.FormValue("grant_type"))
		s.Log.WarnContext(r.Context(), msg)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unsupported_grant_type"})
	}
}

func constructURIWithQueries(uri string, queries map[string]string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("failed url.Parse: %w", err)
	}
	q := u.Query()
	for k, v := range queries {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func parseBasicAuth(auth string) (string, string, error) {
	if !strings.HasPrefix(strings.ToLower(auth), "basic ") {
		return "", "", fmt.Errorf("auth header is not basic: %s", auth)
	}
	decodedAuthContent, err := base64.StdEncoding.DecodeString(auth[len("basic "):])
	if err != nil {
		return "", "", fmt.Errorf("failed base64.StdEncoding.DecodeString: %w", err)
	}
	log.Printf("decoded: %v, %s\n", string(decodedAuthContent), auth)
	clientCredentials := strings.Split(string(decodedAuthContent), ":")
	if len(clientCredentials) != 2 {
		return "", "", fmt.Errorf("basic auth must have two parts: %v", clientCredentials)
	}

	return clientCredentials[0], clientCredentials[1], nil
}
