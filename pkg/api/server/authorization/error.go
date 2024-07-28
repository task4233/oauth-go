package authorization

import (
	"log/slog"
	"net/http"
	"net/url"
)

type ErrorType string

const (
	InvalidRequest ErrorType = "invalid_request"
	ServerError    ErrorType = "server_error"
)

type ErrorResponse struct {
	Error       ErrorType // required
	Description string    // optional
	ErrorURI    string    // optional
	State       string    // optional
}

func RequestError(w http.ResponseWriter, r *http.Request, errResp *ErrorResponse) {
	u := url.Values{}
	u.Set("error", string(errResp.Error))
	u.Set("description", errResp.Description)
	u.Set("state", errResp.State)

	slog.Info("error", slog.String("u", u.Encode()))

	http.Redirect(w, r, errResp.ErrorURI+"?"+u.Encode(), http.StatusFound)
}
