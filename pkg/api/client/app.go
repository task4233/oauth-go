package client

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var (
	//go:embed templates/index.html.tmpl
	indexTemplate string

	//go:embed templates/login.html.tmpl
	loginTemplate string

	//go:embed templates/callback.html.tmpl
	callbackTemplate string
)

type App struct {
	oauthConfig  *oauth2.Config
	stateStorage map[string]struct{}
}

func NewApp(oauthConfig *oauth2.Config) *App {
	return &App{
		oauthConfig:  oauthConfig,
		stateStorage: make(map[string]struct{}),
	}
}

func (s *App) Run(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Index)
	mux.HandleFunc("/login", s.Login)
	mux.HandleFunc("/auth/callback", s.Callback)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (s *App) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexTemplate))
}

func (s *App) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.LoginGET(w, r)
	case http.MethodPost:
		s.LoginPOST(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *App) LoginGET(w http.ResponseWriter, r *http.Request) {
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

func (s *App) LoginPOST(w http.ResponseWriter, r *http.Request) {
	state := uuid.NewString()
	// TODO: need to consider how to recognize each state
	s.stateStorage[state] = struct{}{}

	authCodeURL := s.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, authCodeURL, http.StatusFound)
}

type CallbackRequest struct {
	Code  string // required
	State string // required if an AuthZRequest has the state
}

func (r *CallbackRequest) Validate() error {
	return nil
}

func (s *App) Callback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if params.Get("error") != "" {
		http.Error(w, params.Get("error")+":"+params.Get("error_description"), http.StatusInternalServerError)
		return
	}

	state := params.Get("state")
	if _, ok := s.stateStorage[state]; !ok {
		http.Error(w, "invalid state", http.StatusUnauthorized)
		return
	}

	code := params.Get("code")
	opts := make([]CodeExchangeOption, 0)
	opts = append(opts, func() []oauth2.AuthCodeOption {
		return []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("client_id", s.oauthConfig.ClientID),
		}
	})
	tokens, err := s.CodeExchange(r.Context(), code, opts...)
	if err != nil {
		http.Error(w, "failed to exchange token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	t, err := template.New("callback").Parse(callbackTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, tokens.AccessToken)
	if err != nil {
		slog.Error("failed to execute template", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type CodeExchangeOption func() []oauth2.AuthCodeOption

type CodeExchangeResponse struct {
	AccessToken *oauth2.Token `json:"access_token"`
}

func (s *App) CodeExchange(ctx context.Context, code string, opts ...CodeExchangeOption) (*CodeExchangeResponse, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, http.DefaultClient)

	codeOpts := make([]oauth2.AuthCodeOption, 0)
	for _, opt := range opts {
		codeOpts = append(codeOpts, opt()...)
	}

	token, err := s.oauthConfig.Exchange(ctx, code, codeOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	return &CodeExchangeResponse{AccessToken: token}, nil
}
