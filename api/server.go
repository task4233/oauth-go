package api

import (
	"net/http"

	"github.com/task4233/oauth/logger"
)

type Server interface {
	Run() error
}

// logAdapter is a middleware for common logging for handlers.
func logAdapter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("[Req] %s %s\n", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Info("[Res] %s %s\n", r.Method, r.URL.Path)
	})
}
