package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	introspectTimeout  = 5 * time.Second
	introspectEndpoint = "http://localhost:9001/introspect"
)

func AuthAdapter(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		slog.Info("auth adapter", slog.String("header", fmt.Sprintf("%#v", r.Header)))
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slog.Info("auth header", slog.String("authHeader", authHeader))

		ok, err := verifyToken(extractToken(authHeader))
		if err != nil || !ok {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
			return
		}

		slog.Info("auth ok")

		next.ServeHTTP(w, r)
	})
}

func extractToken(authHeader string) string {
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// ref: https://datatracker.ietf.org/doc/html/rfc7662#section-2.2
type IntrospectResponse struct {
	Active    bool   `json:"active"`     // required
	Scope     string `json:"scope"`      // optional
	ClientID  string `json:"client_id"`  // optional
	Username  string `json:"username"`   // optional
	TokenType string `json:"token_type"` // optional
	Exp       int64  `json:"exp"`        // optional, expiration time
	Iat       int64  `json:"iat"`        // optional, issued at
	Nbf       int64  `json:"nbf"`        // optional, not to be used before
	Sub       string `json:"sub"`        // optional, subject of the token
	Aud       string `json:"aud"`        // optional, audience
	Jti       string `json:"jti"`        // optional, JWT ID
}

func verifyToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), introspectTimeout)
	defer cancel()

	u := &url.Values{}
	u.Set("token", token)
	u.Set("token_type_hint", "access_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, introspectEndpoint, strings.NewReader(u.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ir := &IntrospectResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ir)
	if err != nil {
		return false, err
	}

	return ir.Active, nil
}
