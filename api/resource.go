package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/task4233/oauth-go/logger"
	"golang.org/x/exp/slog"
)

const (
	bearerPrefix   = "bearer "
	accessTokenKey = "access_token"
)

type Resource struct {
	port int
	Log  *slog.Logger
}

func NewResource(ctx context.Context, port int) *Resource {
	return &Resource{
		port: port,
		Log:  logger.FromContext(ctx),
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
	s.Log.InfoContext(r.Context(), "POST / is called")

	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	s.log.Printf("failed io.ReadAll: %s\n", err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// defer r.Body.Close()

	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		r.ParseForm()
	}

	// get access token
	auth := r.Header.Get("authorization")
	inToken := ""
	if len(auth) > 0 && strings.Index(strings.ToLower(auth), bearerPrefix) == 0 {
		inToken = auth[len(bearerPrefix):]
	} else if len(r.Form.Get(accessTokenKey)) > 0 {
		inToken = r.Form.Get(accessTokenKey)
	} else if len(r.URL.Query().Get(accessTokenKey)) > 0 {
		inToken = r.URL.Query().Get(accessTokenKey)
	}

	_ = inToken
	panic("not implemented")
}

type NoSQL interface {
	Get(tableName string, keyName string) (value string)
}
