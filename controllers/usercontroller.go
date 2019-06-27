package controllers

import (
	"net/http"
)

type UserController struct{}

func (c UserController) registerRoutes() {
	http.HandleFunc("/api/v1/user", checkJwt(c.handleUser))
}

//just returns token info
func (c UserController) handleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user"))
}
