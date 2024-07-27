package authentication

import (
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
)

type Authentication struct{}

func NewAuthentication() *Authentication {
	return &Authentication{}
}

func (s *Authentication) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", s.Login)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

//go:embed templates/login.html.tmpl
var loginTemplate string

func (s *Authentication) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.LoginGET(w, r)
	case http.MethodPost:
		s.LoginPOST(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Authentication) LoginGET(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("login").Parse(loginTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Authentication) LoginPOST(w http.ResponseWriter, r *http.Request) {
	client := r.FormValue("client")
	http.Redirect(w, r, "http://localhost:9001/authorize?client="+client, http.StatusFound)
}
