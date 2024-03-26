package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/task4233/oauth/common"
	"github.com/task4233/oauth/domain"
	"github.com/task4233/oauth/infra"
	"github.com/task4233/oauth/infra/repository"
)

var ErrInvalidScope = errors.New("invalid scope")

const expiresDuration = time.Minute

type AuthorizationServer struct {
	srv              http.Server
	clientRepo       repository.Client
	accessTokenRepo  repository.AccessToken
	authReqRepo      repository.AuthorizationRequest
	authCodeRepo     repository.AuthorizationCode
	refreshTokenRepo repository.RefreshToken
	log              *slog.Logger
}

func NewAuthorizationServer(port int, clientRepo repository.Client, accessTokenRepo repository.AccessToken, log *slog.Logger) *AuthorizationServer {
	return &AuthorizationServer{
		srv: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: (&AuthorizationServer{}).route(),
		},
		clientRepo:       clientRepo,
		accessTokenRepo:  accessTokenRepo,
		authReqRepo:      infra.NewAuthorizationRequestRepository(),
		refreshTokenRepo: infra.NewRefreshTokenRepository(),
		log:              log,
	}
}

func (s *AuthorizationServer) route() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/authorize", s.authorize)
	mux.HandleFunc("/approve", s.approve)
	mux.HandleFunc("/token", s.token)

	return mux
}

func (s *AuthorizationServer) Run() error {
	return http.ListenAndServe(s.srv.Addr, s.srv.Handler)
}

func (s *AuthorizationServer) authorize(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// get a client_id from the query
	clientID := r.URL.Query().Get("client_id")

	// find the client from a repository by the client_id
	// if not found, return invalid argument
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a scope from the query
	scope := r.URL.Query().Get("scope")

	// check if all the scope are contained in the client's scope
	// if not return invalid argument
	rScopes := strings.Split(scope, " ")
	for _, scope := range rScopes {
		if !slices.Contains(client.Scopes, scope) {
			http.Error(w, ErrInvalidScope.Error(), http.StatusBadRequest)
			return
		}
	}

	redirectURI := r.URL.Query().Get("redirect_uri")

	// generate reqID
	// save authorization_request in the repository, redirect_uri
	reqID := uuid.NewString()
	state := r.URL.Query().Get("state")
	err = s.authReqRepo.Insert(ctx, &domain.AuthorizationRequest{
		ID:          reqID,
		ClientID:    clientID,
		State:       state,
		Scopes:      rScopes,
		RedirectURI: redirectURI,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to redirect_uri with client_name, reqID, req.scope
	if redirectURI == "" {
		redirectURI = r.RemoteAddr
	}
	redirectURI, err = common.ConstructURLWithQueries(redirectURI, map[string]string{
		"client_name": client.Name,
		"req_id":      reqID,
		"scope":       scope,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func (s *AuthorizationServer) approve(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// get a req_id from the query
	// if request does not exist, return invalid request
	reqID := r.URL.Query().Get("req_id")
	authReq, err := s.authReqRepo.Get(ctx, reqID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if response type is "code"
	// If not, return invalid request
	responseType := r.URL.Query().Get("response_type")
	if responseType != "code" {
		http.Error(w, fmt.Sprintf("response_type: %s is invalid", responseType), http.StatusBadRequest)
		return
	}

	// get user_id from the query
	// TODO: get a session_id from the Cookie
	// if not found, redirect to login page
	userID := r.URL.Query().Get("user_id")

	// get a client_id from the query
	// find the client from a repository by the client_id
	// if not found, return invalid argument
	clientID := r.URL.Query().Get("client_id")
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a scope from the query
	// check if all the scope are contained in the client's scope
	// if not return invalid argument
	scope := r.URL.Query().Get("scope")
	rScopes := strings.Split(scope, " ")
	for _, scope := range rScopes {
		if !slices.Contains(client.Scopes, scope) {
			http.Error(w, ErrInvalidScope.Error(), http.StatusBadRequest)
			return
		}
	}

	redirectURI := r.URL.Query().Get("redirect_uri")
	if authReq.RedirectURI != "" && authReq.RedirectURI != redirectURI {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}

	// generate authorization_code(signature)
	// save authorization code with scope and user
	code := uuid.NewString()
	err = s.authCodeRepo.Insert(ctx, &domain.AuthorizationCode{
		Code:        code,
		UserID:      userID,
		ClientID:    clientID,
		Scopes:      rScopes,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(expiresDuration),
		DisabledAt:  time.Time{},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to the redirect_uri with the code and state
	redirectURI, err = common.ConstructURLWithQueries(redirectURI, map[string]string{
		"code":  code,
		"state": authReq.State,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func (s *AuthorizationServer) token(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// get a client_id and a client_secret from the authorization header
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		http.Error(w, "unauthorized by basic auth", http.StatusUnauthorized)
		return
	}

	// find the client from a repository by the client_id
	// if not found, return invalid argument
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// confirm if client's client_secret and auth.client_secret are same
	// if not, return not authenticated
	clientSecretHashByte := sha256.Sum256([]byte(clientSecret))
	clientSecretHash := hex.EncodeToString(clientSecretHashByte[:])

	if clientSecretHash != client.SecretHash {
		http.Error(w, "invalid client secret", http.StatusBadRequest)
		return
	}

	grantType := r.URL.Query().Get("grant_type")
	switch grantType {
	case "authorization_code":
		//
		// authorization_code grant
		//

		// get a code from the query
		// find the code information from the code
		// if not found, return invalid argument
		code := r.URL.Query().Get("code")
		authCode, err := s.authCodeRepo.Get(ctx, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if the code's client_id and auth.client_id are same
		// if not, return invalid argument
		if authCode.ClientID != clientID {
			http.Error(w, "invalid client_id", http.StatusBadRequest)
			return
		}

		// find the user from code information's user_id
		userID := authCode.UserID

		// generate access_token
		// store access_token with client_id, scope, user_id
		signature := uuid.NewString()
		accessToken := &domain.AccessToken{
			Signature: signature,
			UserID:    userID,
			ClientID:  clientID,
			Scopes:    authCode.Scopes,
			ExpiresAt: time.Now().Add(expiresDuration),
		}
		err = s.accessTokenRepo.Insert(ctx, accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// generate refresh_token
		// store refresh_token with client_id, scope, user_id
		signature = uuid.NewString()
		refreshToken := &domain.RefreshToken{
			Signature:  signature,
			UserID:     userID,
			ClientID:   clientID,
			Scopes:     authCode.Scopes,
			ExpiresAt:  time.Now().Add(expiresDuration),
			DisabledAt: time.Time{},
		}
		err = s.refreshTokenRepo.Insert(ctx, refreshToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return access_token with token_type: bearer, refresh_token, client.scope
		type authorizationCodeResponse struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			RefreshToken string `json:"refresh_token"`
			Scope        string `json:"scope"`
		}
		_ = json.NewEncoder(w).Encode(authorizationCodeResponse{
			AccessToken:  accessToken.Signature,
			TokenType:    "Bearer",
			RefreshToken: refreshToken.Signature,
			Scope:        strings.Join(authCode.Scopes, " "),
		})
		return
		//
		// client_credentials
		//

		// check if grant_type is client_credentaials
		// if not, return invalid argument

		// get a scope from the body
		// check if all the scope are contained in the client's scope
		// if not return invalid argument

		// generate access_token
		// store access_token with client_id, scope, user_id

		// return access_token with token_type: bearer, client.scope

		//
		// refresh_token
		//

		// get a refresh_token from the body
		// find refresh_token info from the repository

		// check if auth.client_id and refresh_token.client_id are same
		// return invalid_refresh_token

		// generate access_token
		// store access_token with client_id, scope, user_id

		// generate refresh_token
		// store refresh_token with client_id, scope, user_id

		// disable the past refresh_token

		// return access_token with token_type: bearer, refresh_token, client.scope
	}

	panic("not implemented")
}
