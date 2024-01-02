package api

type errorResponse int

// ref: https://datatracker.ietf.org/doc/html/rfc6749#section-4.2.2.1
const (
	unknown errorResponse = iota
	invalidRequest
	unauthorizedClient
	accessDenied
	unsupportedResponseType
	invalidScope
	serverError
	temporarilyUnavailable
)

var errorResponseMessages = []string{
	"unknown",
	"invalid_request",
	"unauthorized_client",
	"access_denied",
	"unsupported_response_type",
	"invalid_scope",
	"server_error",
	"temporarily_unavailable",
}

func (e errorResponse) String() string {
	if e < unknown || len(errorResponseMessages) <= int(e) {
		return errorResponseMessages[unknown]
	}
	return errorResponseMessages[e]
}
