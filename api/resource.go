package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/task4233/oauth/infra/repository"
)

// resourceServer is a server for hosting protected resources, capable of accepting
// and responding to protected resource requests using access tokens.
type resourceServer struct {
	port int
	srv  http.Server
}

func NewResourceServer(ctx context.Context, port int, kvs repository.KVS) Server {
	s := &resourceServer{port: port}
	s.srv.Addr = fmt.Sprintf(":%d", port)
	s.srv.Handler = s.route()
	return s
}

func (s *resourceServer) Run() error {
	return s.srv.ListenAndServe()
}

func (a *resourceServer) route() http.Handler {
	h := http.NewServeMux()

	return h
}
