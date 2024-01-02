package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/task4233/oauth/infra/repository"
)

type AuthorizationServer struct {
	srv  http.Server
	repo repository.KVS
	log  *slog.Logger
}

func NewAuthorizationServer(port int, repo repository.KVS, log *slog.Logger) *AuthorizationServer {
	return &AuthorizationServer{
		srv: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: (&AuthorizationServer{}).route(),
		},
		repo: repo,
		log:  log,
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
	return nil
}

func (s *AuthorizationServer) authorize(w http.ResponseWriter, r *http.Request) {
	// get a client_id from the query

	// find the client from a repository by the client_id
	// if not found, return invalid argument

	// get a scope from the query
	// check if all the scope are contained in the client's scope
	// if not return invalid argument

	// generate request_id
	// save authorization_request in the repository, redirect_uri

	// get redirect_uri

	// redirect to redirect_uri with client info, req_id, req.scope

	panic("not implemented")
}

func (s *AuthorizationServer) approve(w http.ResponseWriter, r *http.Request) {
	// get a req_id from the query
	// if request does not exist, return invalid request

	// check if response type is "code"
	// If not, return invalid request

	// get a user_id from the session_id
	// if not found, redirect to login page

	// get a client_id from the query
	// find the client from a repository by the client_id
	// if not found, return invalid argument

	// get a scope from the query
	// check if all the scope are contained in the client's scope
	// if not return invalid argument

	// generate authorization_code(signature)
	// save authorization code with scope and user

	// redirect to the redirect_uri with the code and state

	panic("not implemented")
}

func (s *AuthorizationServer) token(w http.ResponseWriter, r *http.Request) {
	// get a client_id and a client_secret from the authorization header
	// find the client from a repository by the client_id
	// if not found, return invalid argument

	// confirm if client's client_secret and auth.client_secret are same
	// if not, return not authenticated

	//
	// authorization_code grant
	//

	// check if grant_type is authorization_code
	// if not, return invalid argument

	// get a code from the query
	// find the code information from the code
	// if not found, return invalid argument

	// check if the code's client_id and auth.client_id are same
	// if not, return invalid argument

	// find the user from code information's user_id

	// generate access_token
	// store access_token with client_id, scope, user_id

	// generate refresh_token
	// store refresh_token with client_id, scope, user_id

	// return access_token with token_type: bearer, refresh_token, client.scope

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

	panic("not implemented")
}
