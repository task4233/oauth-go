package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/task4233/oauth-go/api"
)

// authorization server information
var authServer = map[string]string{
	"authorizationEndpoint": "http://localhost:9001/authorize",
	"tokenEndpoint":         "http://localhost:9001/token",
}

// client information
var client = api.Client{
	ClientID:     "oauth-client-id-1",
	ClientSecret: "oauth-client-secret-1",
	RedirectURI:  []string{"http://localhost:9000/callback"},
	Scope:        "read write",
}

const protectedResource = "http://localhost:9002/resource"

var state string

func main() {
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/authorize is called: %#v\n", r)

		state = uuid.New().String()

		// authroize request
		v := url.Values{
			"response_type": []string{"code"},
			"scope":         []string{client.Scope},
			"client_id":     []string{client.ClientID},
			"redirect_uri":  client.RedirectURI,
			"state":         []string{state},
		}
		baseURL := authServer["authorizationEndpoint"]
		uStr := fmt.Sprintf("%s?%s", baseURL, v.Encode())

		log.Printf("redirect: %s\n", uStr)
		http.Redirect(w, r, uStr, http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("/callback is called: %#v\n", r)

		errStr := r.URL.Query().Get("error")
		if errStr != "" {
			log.Printf("error is given in callback handler: %s\n", errStr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		reqState := r.URL.Query().Get("state")
		if state == "" || state != reqState {
			log.Printf("invalid state: %s, expected: %s\n", reqState, state)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req, err := http.NewRequest(http.MethodPost, authServer["tokenEndpoint"], nil)
		if err != nil {
			log.Printf("failed to create request: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		code := r.URL.Query().Get("code")
		formData := url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {code},
			"redirect_uri": {client.RedirectURI[0]},
		}
		req.Body = io.NopCloser(strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", client.ClientID, client.ClientSecret)))))

		tokenRes, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("failed to request token: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("token response: %#v\n", tokenRes)

		if tokenRes.StatusCode != http.StatusOK {
			log.Printf("failed to request token: %v\n", tokenRes.Status)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		dat, err := io.ReadAll(tokenRes.Body)
		if err != nil {
			log.Printf("failed to read token response: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		defer tokenRes.Body.Close()
		log.Printf("access token: %s\n", string(dat))
	})

	http.ListenAndServe(":9000", nil)
}
