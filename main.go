package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/iainkiloh/examplegoapi/configurations"
	"github.com/iainkiloh/examplegoapi/controllers"
	"github.com/iainkiloh/examplegoapi/middleware"
	"github.com/iainkiloh/examplegoapi/queries"
)

func main() {

	//load configuration from file
	loadConfig()

	//setup db conection
	db := connectToDatabase()
	defer db.Close()

	//register all routes
	controllers.Startup()

	//setup the middleware pipeline - logging, compression, DefaultServeMux
	pipeline := middleware.NewLoggerMiddleware(
		middleware.NewGzipMiddleware(http.DefaultServeMux))

	//start listening on port, use the setup middleware
	http.ListenAndServe(":9999", pipeline)

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
