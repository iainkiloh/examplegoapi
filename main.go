package main

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"

	"github.com/iainkiloh/examplegoapi/configurations"
	"github.com/iainkiloh/examplegoapi/controllers"
	"github.com/iainkiloh/examplegoapi/middleware"
	"github.com/iainkiloh/examplegoapi/queries"
)

func main() {

	//load configuration from file
	loadConfig()

	setUpJwtMiddleware()

	//setup db conection
	db := connectToDatabase()
	defer db.Close()

	//register all routes
	controllers.Startup()

	//setup the middleware pipeline - timeout, logging, compression, DefaultServeMux
	requestMiddlewarePipeline := middleware.NewTimeoutMiddleware(
		middleware.NewLoggerMiddleware(
			middleware.NewGzipMiddleware(
				http.DefaultServeMux)))

	//start listening on port, use the setup middleware
	http.ListenAndServe(":9999", requestMiddlewarePipeline)

	//log api startup
	fmt.Println("examplegoapi started and listening on port 9999")

}

func loadConfig() {

	//open config file to read configuration info
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, "secrets/examplegoapi.secrets.json")
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatalln("unable to open configuration file", err)
		panic(err)
	}

	//set the globally available configuration instance
	configuration := new(configurations.Configuration)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln("unable to read configuration file")
		panic(err)
	}
	configurations.SetConfiguration(configuration)

}

func connectToDatabase() *sql.DB {

	//attempt to open db connection
	db, err := sql.Open("postgres", configurations.GetConnectionString())
	if err != nil {
		log.Fatalln("unable to connect to db", err)
		panic(err)
	}
	//ensure it is wroking by ping-ing it
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB connected")

	queries.SetDatabase(db)
	return db
}

func setUpJwtMiddleware() {

	signingKey := configurations.GetTokenSigningKey()
	var err error
	decodedKey, err := b64.StdEncoding.DecodeString(signingKey)
	if err != nil {
		panic(err)
	}

	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//set the token issuer and audience from configuration
			iss, aud := configurations.GetTokenIssuerAndAudience()
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			return []byte(decodedKey), nil
		},

		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	middleware.SetJwtMiddleware(jwtMiddleware)
}
