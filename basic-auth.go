package main

import "net/http"

type BasicAuthHandler struct {
	Username    string
	Password    string
	NextHandler http.Handler
}

func NewBasicAuthHandler(username string, password string, handler http.Handler) *BasicAuthHandler {
	return &BasicAuthHandler{
		Username: username,
		Password: password,
		NextHandler: handler,
	}
}

func (h *BasicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, pass, _ := r.BasicAuth()
	if !h.checkCredentials(user, pass) {
		http.Error(w, "Unauthorized", 401)
		return
	}
	h.NextHandler.ServeHTTP(w, r)
}

func (h *BasicAuthHandler) checkCredentials(username string, password string) bool {
	return username == h.Username && password == h.Password
}
