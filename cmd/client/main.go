package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/task4233/oauth/api"
	"github.com/task4233/oauth/logger"
	"golang.org/x/exp/slog"
)

const (
	timeout                 = 5 * time.Second
	clientServerPort        = 8000
	authorizationServerPort = 8080
	resourceServerPort      = 9090
)

func main() {
	ctx := context.Background()
	log := logger.FromContext(ctx)
	clients := client{
		clientID:     "test_client",
		clientSecret: "test_client_secret",
		scope:        "read write",
	}

	clientServer := NewClientServer(ctx, clientServerPort, clients)

	log.Info("client server is running...", "port", clientServerPort)
	if err := clientServer.Run(); err != nil {
		log.Error("failed to run client server", "error", err)
		return
	}
}

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
	s.srv.Addr = fmt.Sprintf(":%d", port)
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
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Error("failed to send request", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.log.Error("failed to get authorization", "status", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	// redirect to callback
	io.Copy(w, resp.Body)
}

func (s *clientServer) callback(w http.ResponseWriter, r *http.Request) {
	// state check
	gotState := r.URL.Query().Get("state")
	if s.client.state != gotState {
		s.log.Error("invalid state", "want", s.client.state, "got", gotState)
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")

	// get token
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", fmt.Sprintf("http://localhost:%d/", clientServerPort))
	params.Add("client_id", s.client.clientID)
	params.Add("state", gotState)

	targetURL := fmt.Sprintf("http://localhost:%d/token?%s",
		authorizationServerPort,
		params.Encode(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, nil)
	if err != nil {
		s.log.Error("failed to create token issue request", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(s.client.clientID, s.client.clientSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Error("failed to send token issue request", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
