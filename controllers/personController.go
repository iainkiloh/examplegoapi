package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/iainkiloh/examplegoapi/contracts"
	"github.com/iainkiloh/examplegoapi/middleware"
	"github.com/iainkiloh/examplegoapi/queries"
)

type PersonController struct{}

func (c PersonController) registerRoutes() {

	//set up route handlers for requests and put them through the claims pipeline
	//which validates the JWT
	//and then adds custom claims to the request context
	http.HandleFunc("/api/v1/person", middleware.CheckJwt(
		middleware.GetClaimsContext(c.handlePersonRoute)))
	http.HandleFunc("/api/v1/person/", middleware.CheckJwt(
		middleware.GetClaimsContext(c.handlePersonRoute)))
}

func (c PersonController) handlePersonRoute(w http.ResponseWriter, r *http.Request) {

	//test the request timeout middleware
	//time.Sleep(12 * time.Second)

	//check for our user in our request context with key of 'customUser'
	userInContext := r.Context().Value("CustomUser")
	fmt.Println("added and fetched from context in controller: ", userInContext)

	switch r.Method {
	case http.MethodPost:
		handlePersonPost(w, r)
	case http.MethodPut:
		handlePersonPut(w, r)
	case http.MethodGet:
		p := strings.Split(r.URL.Path, "/")
		if len(p) == 5 {
			handlePersonGet(w, r)
		} else {
			handlePersonPaged(w, r)
		}
	case http.MethodOptions:
		w.Write([]byte("GET, POST, PUT, PATCH, OPTIONS"))
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handlePersonPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		dec := json.NewDecoder(r.Body)
		var person contracts.PersonForCreate
		err := dec.Decode(&person)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			//save to db and get the id back
			id, err := queries.CreatePersonQuery(&person)
			if err != nil {
				log.Println("error:", err)
			}
			if id == 0 {
				w.WriteHeader(http.StatusInternalServerError)
			}

			//create the CreatedAtUrl and set it in ResponseHeader
			p := strings.Split(r.URL.Path, "/")
			createdAtUrl := "http://" + r.Host + p[0] + "/" + p[1] + "/" + p[2] + "/" + p[3] + "/" + strconv.Itoa(id)
			w.Header().Set("CreatedAt", createdAtUrl)
			w.WriteHeader(http.StatusCreated)

		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handlePersonPut(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		dec := json.NewDecoder(r.Body)
		var person contracts.PersonForUpdate
		err := dec.Decode(&person)
		if err != nil {
			log.Println("error:", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			dbResult, err := queries.UpdatePersonQuery(&person)
			if err != nil {
				log.Println("error:", err)
			}
			if dbResult == 1 {
				//create the UpdatedAt ResponseHeader
				p := strings.Split(r.URL.Path, "/")
				updatedAtUrl := "http://" + r.Host + p[0] + "/" + p[1] + "/" + p[2] + "/" + p[3] + "/" + strconv.Itoa(person.Id)
				w.Header().Set("UpdatedAt", updatedAtUrl)
				w.WriteHeader(http.StatusOK)
			} else if dbResult == 0 {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handlePersonGet(w http.ResponseWriter, r *http.Request) {
	//get the url path segments - note should only be 1
	p := strings.Split(r.URL.Path, "/")
	//check the id converts to an int
	id, err := strconv.Atoi(p[4])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//go to db to fetch
	dbResult, err := queries.GetPersonQuery(id)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dbResult.Id != 0 {
		//encode to Json
		response, err := json.Marshal(dbResult)
		if err != nil {
			log.Println("Error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(string(response)))
	} else {
		w.WriteHeader(http.StatusNotFound) //not found http status code
	}
}

func handlePersonPaged(w http.ResponseWriter, r *http.Request) {
	//check if we have query string params
	params := r.URL.Query() // returns slice of query string params as strings
	orderBy := params.Get("orderBy")
	pageNumber := params.Get("pageNumber")
	itemsPerPage := params.Get("itemsPerPage")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Go - Person Get with route: " +
		" and query params - orderBy:" + orderBy + ", pageNumber: " +
		pageNumber + ", itemsPerPage: " + itemsPerPage))

}
