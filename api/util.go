package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const authCookieName = "auth-token"

func errorResponse(w http.ResponseWriter, err string, status int) {
	http.Error(w, err, status)
}

func okResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func emptyResponse(w http.ResponseWriter, status int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}

func (s *API) setAuthCookie(w http.ResponseWriter, duration time.Duration, token string) {
	cookie := &http.Cookie{
		Name:  authCookieName,
		Value: token,
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}

func (api *API) deleteAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    authCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
}
