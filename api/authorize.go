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

func (s *authorizationServer) Run() error {
	panic("implement me")
}

func (s *authorizationServer) route() http.Handler {
	h := http.NewServeMux()

	h.Handle("/authorize", logAdapter(http.HandlerFunc(s.authorize)))
	h.Handle("/authenticate", logAdapter(http.HandlerFunc(s.authenticate)))
	h.Handle("/token", logAdapter(http.HandlerFunc(s.token)))

	return h
}

// authorize is for handing 2.send the authorization.
func (s *authorizationServer) authorize(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// authenticate is for handling 6.send the information for user authentication.
func (s *authorizationServer) authenticate(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// authenticate is for handling 8.send a token issue request.
func (s *authorizationServer) token(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
