package main

import (
	"encoding/json"
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
	"approveEndpoint":       "http://localhost:9001/approve",
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
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/approve", approveHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/resource", resourceHandler)

	http.ListenAndServe(":9000", nil)
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
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
}

func approveHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/approve is called: %#v\n", r)

	u, err := url.Parse(authServer["approveEndpoint"])
	if err != nil {
		log.Printf("failed to parse approve endpoint: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := u.Query()
	q.Set("req_id", r.URL.Query().Get("req_id"))
	q.Set("redirect_uri", r.URL.Query().Get("redirect_uri"))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		log.Printf("failed to create request: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	formData := url.Values{
		"response_type": {"code"},
		"scope":         {r.URL.Query().Get("scope")},
		"client_id":     {r.URL.Query().Get("client_id")},
		"state":         {r.URL.Query().Get("state")},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {client.RedirectURI[0]},
	}
	req.Body = io.NopCloser(strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client.ClientID, client.ClientSecret)
	log.Printf("req header: %v\n", req.Header)

	acceptRes, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to request token: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("token response: %#v\n", acceptRes)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
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
	req.SetBasicAuth(client.ClientID, client.ClientSecret)

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
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/resource is called: %#v\n", r)
	req, err := http.NewRequest(http.MethodPost, protectedResource, nil)
	if err != nil {
		msg := fmt.Sprintf("failed to create request: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, msg)
		return
	}

	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		msg := "no access token is given"
		log.Println(msg)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, msg)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		msg := fmt.Sprintf("failed to request resource: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, msg)
		return
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read resource response: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, msg)
		return
	}
	defer res.Body.Close()

	log.Printf("resource response: %s\n", string(dat))
	json.NewEncoder(w).Encode(string(dat))
}

/*
http://localhost:9000/approve?client_id=oauth-client-id-1&redirect_uri=http%3A%2F%2Flocalhost%3A9000%2Fcallback&response_type=code&scope=read+write&state=a58f36e4-7d67-4a5d-a529-ca5f866a9044&req_id=2a131d2b-0b3f-4a85-90db-6c719fddb776
*/
