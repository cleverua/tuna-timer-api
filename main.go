package main

import (
	"net/http"
	"os"

	"log"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/justinas/alice"
	"github.com/pavlo/slack-time/utils"
	"github.com/pavlo/slack-time/web"
)

const version = "0.1.0"

var environment *utils.Environment
var connection *gorm.DB

func main() {
	environment, connection = utils.NewEnvironment(getEnvironmentName(), version)
	utils.PrintBanner(environment)

	environment.MigrateDatabase(connection.DB()) //todo: check config option or env variable before doing this

	handlers := web.NewHandlers(environment, connection)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.Health).Methods("GET")
	router.HandleFunc("/health", handlers.Health).Methods("GET")
	router.HandleFunc("/api/v1/timer", handlers.Timer).Methods("POST")

	// Slack will sometimes call the API method using a GET request
	// to check SSL certificate - so we reply with a status handler here
	router.HandleFunc("/api/v1/timer", handlers.Health).Methods("GET")

	defaultMiddleware := alice.New(
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	)

	log.Fatal(http.ListenAndServe(":8080", defaultMiddleware.Then(router)))
}

func getEnvironmentName() string {
	env := os.Getenv("SLACK_TIME_ENV")
	if env == "" {
		env = utils.DevelopmentEnv
	}
	return env
}
