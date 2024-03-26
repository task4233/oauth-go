package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/task4233/oauth/infra/repository"
)

type ResourceServer struct {
	srv  http.Server
	repo repository.KVS
	log  *slog.Logger
}

func NewResourceServer(port int, repo repository.KVS, log *slog.Logger) *ResourceServer {
	return &ResourceServer{
		srv: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: (&ResourceServer{}).route(),
		},
		repo: repo,
		log:  log,
	}
}

func (s *ResourceServer) route() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.index)

	return mux
}

func (s *ResourceServer) Run() error {
	return http.ListenAndServe(s.srv.Addr, s.srv.Handler)
}

func (s *ResourceServer) index(w http.ResponseWriter, r *http.Request) {

}
