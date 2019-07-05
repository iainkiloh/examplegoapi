package controllers

import (
	"net/http"
	"strings"
)

var personController PersonController
var healthController HealthController
var userController UserController

//registers routes
func Startup() {
	personController.registerRoutes()
	healthController.registerRoutes()
	userController.registerRoutes()
}

func SetAvailableAtRouteHeader(w http.ResponseWriter, r *http.Request, id string, created bool) {

	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}
	p := strings.Split(r.URL.Path, "/")
	availabilityUrl := scheme + "://" + r.Host + p[0] + "/" + p[1] + "/" + p[2] + "/" + p[3] + "/" + id
	if created {
		w.Header().Set("CreatedAt", availabilityUrl)
	} else {
		w.Header().Set("UpdatedAt", availabilityUrl)
	}
}
