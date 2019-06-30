package controllers

import (
	"net/http"

	"github.com/iainkiloh/examplegoapi/middleware"
	"github.com/iainkiloh/examplegoapi/queries"
)

type HealthController struct{}

func (c HealthController) registerRoutes() {
	http.HandleFunc("/api/v1/health", middleware.CheckJwt(c.handleHealth))
	http.HandleFunc("/api/v1/health/live", middleware.CheckJwt(c.handleHealthLive))
}

//just returns ok
func (c HealthController) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

//checks db connectivity and returns
func (c HealthController) handleHealthLive(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := queries.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy and live"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
