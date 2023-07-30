package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/task4233/oauth-go/infra/repository"
	"github.com/task4233/oauth-go/logger"
	"golang.org/x/exp/slog"
)

type contextKeyType string

const (
	bearerPrefix                  = "bearer "
	accessTokenKey contextKeyType = "access_token"
)

type Resource struct {
	port         int
	Log          *slog.Logger
	KVS          repository.KVS
	resourceData map[string]string
}

func NewResource(ctx context.Context, port int, kvs repository.KVS, resourceData map[string]string) *Resource {
	return &Resource{
		port:         port,
		Log:          logger.FromContext(ctx),
		KVS:          kvs,
		resourceData: resourceData,
	}
}

func (s *Resource) Run(ctx context.Context) error {
	http.HandleFunc("/", s.resource)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *Resource) resource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.resoruceGet(w, r)
		return
	case http.MethodPost:
		s.resourcePost(w, r)
		return
	default:
		s.Log.WarnContext(r.Context(), "method: %s is not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *Resource) resoruceGet(w http.ResponseWriter, r *http.Request) {
	s.Log.InfoContext(r.Context(), "GET / is called")

	panic("not implemented")
}

func (s *Resource) resourcePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.Log.InfoContext(ctx, "POST / is called")

	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		r.ParseForm()
	}

	// get access token
	auth := r.Header.Get("authorization")
	inToken := ""
	if len(auth) > 0 && strings.Index(strings.ToLower(auth), bearerPrefix) == 0 {
		inToken = auth[len(bearerPrefix):]
	} else if len(r.Form.Get(string(accessTokenKey))) > 0 {
		inToken = r.Form.Get(string(accessTokenKey))
	} else if len(r.URL.Query().Get(string(accessTokenKey))) > 0 {
		inToken = r.URL.Query().Get(string(accessTokenKey))
	}

	tokens, err := s.KVS.Get(inToken)
	if err != nil {
		msg := fmt.Sprintf("failed to find token with %s: %s", inToken, err.Error())
		s.Log.WarnContext(ctx, msg)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, msg)
		return
	}

	ctx = context.WithValue(ctx, accessTokenKey, tokens["access_token"])

	accessToken := ctx.Value(accessTokenKey).(string)
	if accessToken == "" {
		s.Log.WarnContext(ctx, "failed to get access_token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.resourceData)
}
