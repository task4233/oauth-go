package api

import (
	"fmt"
	"net/http"
)

// authorizationServer is a server for issueing access tokens to the client
// after successfully authenticating the resource owner and obtaining authorization.
type authorizationServer struct {
	port int
	srv  http.Server
}

func NewAuthorizationServer(port int) Server {
	s := &authorizationServer{port: port}
	s.srv.Addr = fmt.Sprintf(":%d", port)
	s.srv.Handler = s.route()
	return s
}

func (a *authorizationServer) Run() error {
	panic("implement me")
}

func (a *authorizationServer) route() http.Handler {
	panic("implement me")
}
