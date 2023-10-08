package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/task4233/oauth/api"
	"github.com/task4233/oauth/logger"
	"golang.org/x/exp/slog"
)

const (
	timeout = 5 * time.Second
)

type client struct {
	clientID     string
	clientSecret string
	scope        string
	state        string
}

type clientServer struct {
	port   int
	srv    http.Server
	client client
	log    *slog.Logger
}

func NewClientServer(ctx context.Context, port int, c client) *clientServer {
	s := &clientServer{port: port}
	s.srv.Addr = fmt.Sprintf("/:%d", port)
	s.srv.Handler = s.route()
	s.client = c
	s.log = logger.FromContext(ctx)
	return s
}

func (s *clientServer) Run() error {
	return s.srv.ListenAndServe()
}

func (s *clientServer) route() http.Handler {
	h := http.NewServeMux()

	h.Handle("/authorize", api.LogAdapter(http.HandlerFunc(s.authorize)))
	h.Handle("/callback", api.LogAdapter(http.HandlerFunc(s.callback)))

	return h
}

func (s *clientServer) authorize(w http.ResponseWriter, r *http.Request) {
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", clientServerPort)
	s.client.state = uuid.NewString()

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", s.client.clientID)
	params.Add("state", s.client.state)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", s.client.scope)

	targetURL := fmt.Sprintf("http://localhost:%d/authorize?%s",
		authorizationServerPort,
		params.Encode(),
	)

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		s.log.Error("failed to create request", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Error("failed to send request", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	respBody := &api.AuthorizeResponse{}
	_ = json.NewDecoder(resp.Body).Decode(&respBody)

	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func (s *clientServer) callback(w http.ResponseWriter, r *http.Request) {
	// state check
	gotState := r.URL.Query().Get("state")
	if s.client.state != gotState {
		s.log.Error("invalid state", "want", s.client.state, "got", gotState)
		http.Error(w, "invalid state", http.StatusInternalServerError)
		return
	}
}
