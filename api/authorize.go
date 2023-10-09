package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/task4233/oauth/common"
	"github.com/task4233/oauth/infra/repository"
	"github.com/task4233/oauth/logger"
	"golang.org/x/exp/slog"
)

const dummyToken = "dummy_token"

type Client struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

type AuthorizeResponse struct {
	Client Client `json:"client"`
	State  string `json:"state"`
	ReqID  string `json:"req_id"`
	Scope  string `json:"scope"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type code struct {
	req    url.Values
	scope  string
	userID string
}

var _ (Server) = (*authorizationServer)(nil)

// authorizationServer is a server for issueing access tokens to the client
// after successfully authenticating the resource owner and obtaining authorization.
type authorizationServer struct {
	port     int
	srv      http.Server
	kvs      repository.KVS
	clients  map[string]*Client
	requests map[string]url.Values
	codes    map[string]*code
	log      *slog.Logger
}

func NewAuthorizationServer(
	ctx context.Context,
	port int,
	clients map[string]*Client,
	kvs repository.KVS,
) *authorizationServer {
	s := &authorizationServer{port: port}
	s.srv.Addr = fmt.Sprintf(":%d", port)
	s.srv.Handler = s.route()
	s.kvs = kvs
	s.clients = clients
	s.requests = make(map[string]url.Values)
	s.codes = make(map[string]*code)
	s.log = logger.FromContext(ctx)
	return s
}

func (s *authorizationServer) Run() error {
	return s.srv.ListenAndServe()
}

func (s *authorizationServer) route() http.Handler {
	h := http.NewServeMux()

	h.Handle("/authorize", LogAdapter(http.HandlerFunc(s.authorize)))
	h.Handle("/approve", LogAdapter(http.HandlerFunc(s.approve)))
	h.Handle("/authenticate", LogAdapter(http.HandlerFunc(s.authenticate)))
	h.Handle("/token", LogAdapter(http.HandlerFunc(s.token)))

	return h
}

// authorize is for handing 2.send the authorization.
func (s *authorizationServer) authorize(w http.ResponseWriter, r *http.Request) {
	// get query parameters
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	state := r.URL.Query().Get("state")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")

	// if there's no correct authorization header, redirect to authenticate.
	authorization := r.Header.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer "+dummyToken+":") {
		s.log.Error("/authorize", "msg", "failed to authenticate", "invalid token", authorization)
		http.Redirect(w, r, "/authenticate", http.StatusFound)
		return
	}
	userID := strings.TrimPrefix(authorization, "Bearer "+dummyToken+":")

	// validate query parameters
	if responseType != "code" {
		s.log.Error("/authorize", "invalid response_type", responseType)
		s.handleError(w, r, map[string]string{
			"error": unsupportedResponseType.String(),
			"state": state,
		})
		return
	}

	client, ok := s.clients[clientID]
	if !ok {
		s.log.Error("/authorize", "invalid client_id", clientID)
		s.handleError(w, r, map[string]string{
			"error": unauthorizedClient.String(),
			"state": state,
		})
		return
	}

	// check scope
	err := s.isValidScope(scope, client)
	if err != nil {
		s.log.Error("/authorize", "invalid scope", scope, "error", err)
		s.handleError(w, r, map[string]string{
			"error": invalidScope.String(),
			"state": state,
		})
		return
	}

	// generate req_id and store the request
	reqID := uuid.New().String()
	s.requests[reqID] = r.URL.Query()

	// generate code and store the code
	c := uuid.New().String()
	s.codes[c] = &code{
		req:    s.requests[reqID],
		scope:  scope,
		userID: userID,
	}

	// redirect to redirect_uri
	redirectURI, err = common.ConstructURLWithQueries(redirectURI, map[string]string{
		"code":  c,
		"state": state,
	})
	if err != nil {
		s.log.Error("failed to constructURIWithQueries: %v", err)
		s.handleError(w, r, map[string]string{
			"error": serverError.String(),
			"state": state,
		})
		return
	}

	http.Redirect(w, r, redirectURI, http.StatusFound)
}

// approve is for handling 6.send the information for user authentication.
// this method is not defined in the RFC.
func (s *authorizationServer) approve(w http.ResponseWriter, r *http.Request) {
	// get query parameters
	reqID := r.URL.Query().Get("req_id")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")
	userID := r.URL.Query().Get("user_id")
	redirectURI := r.URL.Query().Get("redirect_uri")

	// validate approve
	req, ok := s.requests[reqID]
	if !ok {
		s.log.Warn("/approve", "invalid req_id", reqID)
		s.handleError(w, r, map[string]string{
			"error": invalidRequest.String(),
		})
		return
	}
	if scope == "" {
		s.log.Warn("/approve", "msg", "scope is empty")
		s.handleError(w, r, map[string]string{
			"error": invalidRequest.String(),
		})
		return
	}
	if userID == "" {
		s.log.Warn("/approve", "msg", "user_id is empty")
		s.handleError(w, r, map[string]string{
			"error": accessDenied.String(),
		})
		return
	}
	if redirectURI == "" {
		s.log.Warn("/approve", "msg", "redirect_uri is empty")
		s.handleError(w, r, map[string]string{
			"error": invalidRequest.String(),
		})
		return
	}

	// generate code and store the code
	c := uuid.New().String()
	s.codes[c] = &code{
		req:    req,
		scope:  scope,
		userID: userID,
	}

	// redirect to redirect_uri
	redirectURI, err := common.ConstructURLWithQueries(redirectURI, map[string]string{
		"code":  c,
		"state": state,
	})
	if err != nil {
		s.log.Error("/approve", "constructURIWithQueries", err)
		s.handleError(w, r, map[string]string{
			"error": serverError.String(),
			"state": state,
		})
		return
	}

	http.Redirect(w, r, redirectURI, http.StatusFound)
}

// token is for handling 8.send a token issue request.
func (s *authorizationServer) token(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	clientID, clientSecret, err := parseBasicAuth(auth)
	if err != nil {
		s.log.Error("/token", "failed to parseBasicAuth", err)
		s.handleError(w, r, map[string]string{
			"error": invalidRequest.String(),
		})
		return
	}

	// get query parameters
	grantType := r.URL.Query().Get("grant_type")
	code := r.URL.Query().Get("code")
	redirectURL := r.URL.Query().Get("redirect_uri")
	cID := r.URL.Query().Get("client_id")

	// validate query parameters
	if grantType != "authorization_code" {
		s.log.Error("/token", "invalid grant_type", grantType)
		s.handleError(w, r, map[string]string{
			"error": unsupportedResponseType.String(),
		})
		return
	}
	c, ok := s.codes[code]
	if !ok {
		s.log.Error("/token", "invalid code", code)
		s.handleError(w, r, map[string]string{
			"error": invalidRequest.String(),
		})
		return
	}
	if redirectURL == "" {
		s.log.Error("/token", "invalid redirect_uri", redirectURL)
		s.handleError(w, r, map[string]string{
			"error": invalidRequest.String(),
		})
		return
	}
	client, ok := s.clients[cID]
	if !ok {
		s.log.Error("/token", "invalid client_id", clientID)
		s.handleError(w, r, map[string]string{
			"error": unauthorizedClient.String(),
		})
		return
	}

	// validate client credentials
	if client.ClientID != clientID || client.ClientSecret != clientSecret {
		s.log.Error("/token", "msg", "invalid client credentials", "clientID", clientID, "clientSecret", clientSecret)
		s.handleError(w, r, map[string]string{
			"error": unauthorizedClient.String(),
		})
		return
	}

	// TODO: not to use uuid
	accessToken := uuid.New().String()
	vv := map[string]string{
		"access_token": accessToken,
		"client_id":    clientID,
		"scope":        c.scope,
	}
	err = s.kvs.Set(accessToken, vv)
	if err != nil {
		s.log.Error("/token", "msg", "failed to Set", "error", err)
		s.handleError(w, r, map[string]string{
			"error": serverError.String(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		Scope:       c.scope,
	})
}

func (s *authorizationServer) handleError(w http.ResponseWriter, r *http.Request, queryParameters map[string]string) {
	redirectURI, err := common.ConstructURLWithQueries(r.URL.Query().Get("redirect_uri"), queryParameters)
	if err != nil {
		msg := fmt.Sprintf("failed to constructURIWithQueries: %v", err)
		s.log.Warn(msg)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func (s *authorizationServer) isValidScope(reqScope string, client *Client) error {
	reqScopes := strings.Split(reqScope, " ")
	clientScopes := strings.Split(client.Scope, " ") // it can be cached.
	if !common.AreTwoUnorderedSlicesSame(reqScopes, clientScopes) {
		return fmt.Errorf("invalid scope, want: %v, req: %v", client.Scope, reqScope)
	}
	return nil
}
