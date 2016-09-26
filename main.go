package main

import (
	"net/http"
	"os"

	"log"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/tuna-timer/tuna-timer-api/utils"
	"github.com/tuna-timer/tuna-timer-api/web"
	"time"
)

const version = "0.1.0"

func main() {

	time.Local = time.UTC

	environment := utils.NewEnvironment(getEnvironmentName(), version)
	utils.PrintBanner(environment)

	session, err := utils.ConnectToDatabase(environment.Config)
	if err != nil {
		log.Fatalf("Failed to connect to Database: %s", err)
	}

	environment.MigrateDatabase(session)
	handlers := web.NewHandlers(environment, session)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.Health).Methods("GET")
	router.HandleFunc("/health", handlers.Health).Methods("GET")
	router.HandleFunc("/api/v1/timer", handlers.Timer).Methods("POST")

	// Slack will sometimes call the API method using a GET request
	// to check SSL certificate - so we reply with a status handler here
	router.HandleFunc("/api/v1/timer", handlers.Health).Methods("GET")
	router.HandleFunc("/temporary/clear_data", handlers.ClearAllData).Methods("GET")

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
