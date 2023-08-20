package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/task4233/oauth/infra"
	"github.com/task4233/oauth/infra/mock"
	"go.uber.org/mock/gomock"
)

var testClients = map[string]*Client{
	"client_id": {
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		Scope:        "read write",
	},
}

func TestAuthorize(t *testing.T) {
	t.Parallel()

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
				Client: *testClients["client_id"],
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
			s := NewAuthorizationServer(context.Background(), 8080, testClients, infra.NewKVS())

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

func TestApprove(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		reqParams map[string]string
	}{
		"ok": {
			reqParams: map[string]string{
				"req_id":       uuid.NewString(),
				"scope":        "read write",
				"state":        uuid.NewString(),
				"user_id":      uuid.NewString(),
				"redirect_uri": "http://localhost/redirect_uri",
			},
		},
		"ng:empty req id": {
			reqParams: map[string]string{
				"req_id": "",
			},
		},
		"ng:empty scope": {
			reqParams: map[string]string{
				"req_id": uuid.NewString(),
				"scope":  "",
			},
		},
		"ng:empty user id": {
			reqParams: map[string]string{
				"req_id":  uuid.NewString(),
				"scope":   "read write",
				"state":   uuid.NewString(),
				"user_id": "",
			},
		},
		"ng:empty redirect id": {
			reqParams: map[string]string{
				"req_id":       uuid.NewString(),
				"scope":        "read write",
				"state":        uuid.NewString(),
				"user_id":      uuid.NewString(),
				"redirect_uri": "",
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// readyServer
			s := NewAuthorizationServer(context.Background(), 8080, testClients, infra.NewKVS())

			// ready request/reseponse and call the enpoint.
			req := httptest.NewRequest(http.MethodGet, "http://localhost/authorize", nil)
			q := req.URL.Query()
			for k, v := range tt.reqParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			got := httptest.NewRecorder()
			s.approve(got, req)

			resp := got.Result()
			if resp == nil {
				t.Fatalf("failed to get the result")
			}
			defer resp.Body.Close()

			// check response
			if http.StatusFound != resp.StatusCode {
				t.Fatalf("unexpected statusCode, want: %v, got: %v", http.StatusFound, resp.StatusCode)
			}
		})
	}
}

func TestToken(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setupMocks        func(mock.MockKVS)
		client            *Client
		reqParams         map[string]string
		wantStatusCode    int
		wantTokenResponse TokenResponse
	}{
		"ok": {
			setupMocks: func(mk mock.MockKVS) {
				mk.EXPECT().Set(gomock.Any(), gomock.Any())
			},
			client: testClients["client_id"],
			reqParams: map[string]string{
				"grant_type":   "authorization_code",
				"code":         "dummy-code",
				"redirect_uri": "http://localhost/redirect",
				"client_id":    "client_id",
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// readyServer
			s := NewAuthorizationServer(context.Background(), 8080, testClients, infra.NewKVS())
			s.codes["dummy-code"] = &code{
				req:    url.Values{},
				scope:  "read write",
				userID: "dummy-user-id",
			}

			// ready request/reseponse and call the enpoint.
			req := httptest.NewRequest(http.MethodGet, "http://localhost/authorize", nil)
			q := req.URL.Query()
			for k, v := range tt.reqParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			req.SetBasicAuth(tt.client.ClientID, tt.client.ClientSecret)
			got := httptest.NewRecorder()
			s.token(got, req)

			resp := got.Result()
			if resp == nil {
				t.Fatalf("failed to get the result")
			}
			defer resp.Body.Close()

			// check response
			if tt.wantStatusCode != resp.StatusCode {
				t.Fatalf("unexpected statusCode, want: %v, got: %v", http.StatusFound, resp.StatusCode)
			}

			tr := TokenResponse{}
			if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
				t.Fatalf("falied to decode: %v", err)
			}
			if tr.AccessToken == "" {
				t.Fatalf("access token must be empty")
			}
			if tr.Scope != tt.client.Scope {
				t.Fatalf("unexpected scope")
			}
			if tr.TokenType != "Bearer" {
				t.Fatalf("unexpected token type: %v", tr.TokenType)
			}
		})
	}
}
