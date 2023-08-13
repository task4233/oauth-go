package api

import (
	"fmt"
	"net/http"
)

// resourceServer is a server for hosting protected resources, capable of accepting
// and responding to protected resource requests using access tokens.
type resourceServer struct {
	port int
	srv  http.Server
}

func NewResourceServer(port int) Server {
	s := &resourceServer{port: port}
	s.srv.Addr = fmt.Sprintf(":%d", port)
	s.srv.Handler = s.route()
	return s
}

func (s *resourceServer) Run() error {
	panic("implement me")
}

func (a *resourceServer) route() http.Handler {
	panic("implement me")
}
