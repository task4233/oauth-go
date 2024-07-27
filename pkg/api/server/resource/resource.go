package server

import "net/http"

type Resource struct{}

func NewResource() *Resource {
	return &Resource{}
}

func (s *Resource) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/resource", s.Resource)
}

func (s *Resource) Resource(w http.ResponseWriter, r *http.Request) {
}
