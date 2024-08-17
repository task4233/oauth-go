package resource

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Resource struct{}

func NewResource() *Resource {
	return &Resource{}
}

func (s *Resource) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/resource", AuthAdapter(s.Resource))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (s *Resource) Resource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,cors")

	switch r.Method {
	case http.MethodGet:
		s.ResourceGET(w, r)
	case http.MethodOptions:
		s.ResourceOptions(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Resource) ResourceGET(w http.ResponseWriter, r *http.Request) {
	slog.Info("resource get")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

func (s *Resource) ResourceOptions(w http.ResponseWriter, r *http.Request) {
}
