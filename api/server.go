package api

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/task4233/oauth/logger"
)

type Server interface {
	Run() error
}

// LogAdapter is a middleware for common logging for handlers.
func LogAdapter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("[Req] %s %s\n", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Info("[Res] %s %s\n", r.Method, r.URL.Path)
	})
}

func parseBasicAuth(auth string) (string, string, error) {
	if !strings.HasPrefix(strings.ToLower(auth), "basic ") {
		return "", "", fmt.Errorf("auth header is not basic: %s", auth)
	}
	decodedAuthContent, err := base64.StdEncoding.DecodeString(auth[len("basic "):])
	if err != nil {
		return "", "", fmt.Errorf("failed base64.StdEncoding.DecodeString: %w", err)
	}
	log.Printf("decoded: %v, %s\n", string(decodedAuthContent), auth)
	clientCredentials := strings.Split(string(decodedAuthContent), ":")
	if len(clientCredentials) != 2 {
		return "", "", fmt.Errorf("basic auth must have two parts: %v", clientCredentials)
	}

	return clientCredentials[0], clientCredentials[1], nil
}
