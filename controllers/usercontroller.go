package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/iainkiloh/examplegoapi/middleware"
)

type UserController struct{}

func (c UserController) registerRoutes() {

	//set up route handlers for requests and put them through the claims pipeline
	//which validates the JWT
	//and then adds custom claims to the request context
	http.HandleFunc("/api/v1/user", middleware.CheckJwt(
		middleware.GetClaimsContext(c.handleUser)))

}

//just returns token info
func (c UserController) handleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userInContext := r.Context().Value("CustomUser")
	//fmt.Println("added and fetched from context in controller: ", tester)
	response, err := json.Marshal(userInContext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(response)))
}
