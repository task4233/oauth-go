package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/task4233/oauth/infra"
)

func TestAuthorize(t *testing.T) {
	t.Parallel()

	clients := map[string]*Client{"client_id": {
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		Scope:        "read write",
	},
	}

	tests := map[string]struct {
		clients   map[string]*Client
		reqParams map[string]string
		wantCode  int
		wantResp  *AuthorizeResponse
	}{
		"ok:correctly done authorization request": {
			reqParams: map[string]string{
				"response_type": "code",
				"client_id":     "client_id",
				"state":         "dummy-state",
				"scope":         "read write",
			},
			wantCode: http.StatusOK,
			wantResp: &AuthorizeResponse{
				Client: *clients["client_id"],
				State:  "dummy-state",
				Scope:  "read write",
			},
		},
		"ng:invalid response type": {
			reqParams: map[string]string{
				"response_type": "invalid_code",
			},
			wantCode: http.StatusFound,
		},
		"ng:invalid clientID": {
			reqParams: map[string]string{
				"response_type": "code",
				"client_id":     "invalid_client_id",
			},
			wantCode: http.StatusFound,
		},
		"ng:invalid scopes": {
			reqParams: map[string]string{
				"response_type": "code",
				"client_id":     "client_id",
				"scope":         "invalid_scope",
			},
			wantCode: http.StatusFound,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// readyServer
			s := NewAuthorizationServer(context.Background(), 8080, clients, infra.NewKVS())

			// ready request/reseponse and call the enpoint.
			req := httptest.NewRequest(http.MethodGet, "http://localhost/authorize", nil)
			q := req.URL.Query()
			for k, v := range tt.reqParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			got := httptest.NewRecorder()
			s.authorize(got, req)

			resp := got.Result()
			if resp == nil {
				t.Fatalf("failed to get the result")
			}
			defer resp.Body.Close()

			// check response
			if tt.wantCode != resp.StatusCode {
				t.Fatalf("unexpected statusCode, want: %v, got: %v", tt.wantCode, resp.StatusCode)
			}
			if tt.wantResp != nil {
				respBody := AuthorizeResponse{}
				err := json.NewDecoder(resp.Body).Decode(&respBody)
				if err != nil {
					t.Fatalf("failed to decode json: %v", err)
				}

				if diff := cmp.Diff(tt.wantResp, &respBody, cmpopts.IgnoreFields(AuthorizeResponse{}, "ReqID")); diff != "" {
					t.Fatalf("unexpected response(-want+got): %s", diff)
				}
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
	}{
		"success": {},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_ = tt
		})
	}
}

func TestToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{}{
		"success": {},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_ = tt
		})
	}
}
