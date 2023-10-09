package api

import (
	"fmt"
	"net/http"
)

func (s *authorizationServer) authenticate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("user_id") == "" {
		s.log.Warn("/authenticate", "msg", "failed to authenticate because of empty client_id")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s:%s", dummyToken, r.URL.Query().Get("user_id")))
}
