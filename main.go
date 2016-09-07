package main

import (
	"net/http"

	"log"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/pavlo/slack-time/utils"
	"github.com/pavlo/slack-time/web"
)

const version = "0.0.1"

var status = map[string]string{"version": version}
var environment *utils.Environment

func main() {
	environment = utils.NewEnvironment(utils.DevelopmentEnv, version)
	handlers := web.NewHandlers(environment)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.Health).Methods("GET")
	router.HandleFunc("/health", handlers.Health).Methods("GET")
	router.HandleFunc("/api/v1/timer", handlers.Timer).Methods("POST")

	// Slack will sometimes call the API method with a GET request
	// to check SSL certificate so we reply with a status handler here
	router.HandleFunc("/api/v1/timer", handlers.Health).Methods("GET")

	defaultMiddleware := alice.New(
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	)

	log.Fatal(http.ListenAndServe(":8080", defaultMiddleware.Then(router)))
}
