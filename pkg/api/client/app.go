package client

import "net/http"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (s *App) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /login", s.LoginGET)
	mux.HandleFunc("POST /login", s.LoginPOST)
}

func (s *App) LoginGET(w http.ResponseWriter, r *http.Request) {
	// TODO:
}

func (s *App) LoginPOST(w http.ResponseWriter, r *http.Request) {
}
