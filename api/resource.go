package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/task4233/oauth/infra/repository"
	"github.com/task4233/oauth/logger"
	"golang.org/x/exp/slog"
)

// resourceServer is a server for hosting protected resources, capable of accepting
// and responding to protected resource requests using access tokens.
type resourceServer struct {
	port int
	srv  http.Server
	kvs  repository.KVS
	log  *slog.Logger
}

func NewResourceServer(ctx context.Context, port int, kvs repository.KVS) Server {
	s := &resourceServer{port: port}
	s.srv.Addr = fmt.Sprintf(":%d", port)
	s.srv.Handler = s.route()
	s.kvs = kvs
	s.log = logger.FromContext(ctx)
	return s
}

func (s *resourceServer) Run() error {
	return s.srv.ListenAndServe()
}

func (s *resourceServer) route() http.Handler {
	h := http.NewServeMux()
	h.HandleFunc("/", s.index)
	return h
}

func (s *resourceServer) index(w http.ResponseWriter, r *http.Request) {
	// extract bearer token from Authorization header
	authHeader := strings.ToLower(r.Header.Get("Authorization"))
	inToken, ok := strings.CutPrefix(authHeader, "bearer ")
	if !ok {
		s.log.Error("/index", "msg", "failed to extract bearer token from Authorization header", "Authorization header", authHeader)
		http.Error(w, fmt.Errorf("failed to extract bearer token from Authorization header: %s", authHeader).Error(), http.StatusBadRequest)
		return
	}

	// extract token from kvs
	_, err := s.kvs.Get(inToken)
	if err != nil {
		s.log.Error("/index", "msg", "failed to extract token from kvs", "error", err)
		http.Error(w, fmt.Errorf("failed to extract token from kvs: %w", err).Error(), http.StatusUnauthorized)
		return
	}

	w.Write([]byte("You're authorized!"))
}
