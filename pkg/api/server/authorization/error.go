package authorization

import (
	"net/http"
	"net/url"

	"github.com/task4233/oauth/pkg/domain/model"
)

type errorType string

const (
	InvalidRequest errorType = "invalid_request"
	ServerError    errorType = "server_error"
)

type ErrorResponse struct {
	Error       errorType // required
	Description string    // optional
	ErrorURI    string    // optional
	State       string    // optional
}

func AuthRequestError(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	if authReq == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ErrorResponse{
		Error:       ServerError,
		Description: err.Error(),
		ErrorURI:    authReq.RedirectURI,
		State:       authReq.State,
	}

	u := url.Values{}
	u.Set("error", string(resp.Error))
	u.Set("description", resp.Description)
	u.Set("state", resp.State)
	http.Redirect(w, r, resp.ErrorURI+"?"+u.Encode(), http.StatusFound)
}
