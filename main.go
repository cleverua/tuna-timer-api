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

	"fmt"
	"github.com/robfig/cron"
	"github.com/tuna-timer/tuna-timer-api/jobs"
	"gopkg.in/mgo.v2"
)

const (
	version = "0.1.0"
	port    = "8080"
)

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

	router.HandleFunc("/api/v1/health", handlers.Health).Methods("GET")

	// Slack will sometimes call the API method using a GET request
	// to check SSL certificate - so we reply with a status handler here
	router.HandleFunc("/api/v1/timer", handlers.Timer).Methods("POST", "GET")

	// Slack  OAuth2 stuff
	router.HandleFunc("/api/v1/slack/oauth2redirect", handlers.SlackOauth2Redirect).Methods("GET")

	// Static assets
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	// Temporary stuff, remove eventually
	router.HandleFunc("/api/v1/temporary/clear_data", handlers.ClearAllData).Methods("GET")
	router.HandleFunc("/api/v1/temporary/send_message", handlers.SendSampleMessageFromBot).Methods("GET")

	defaultMiddleware := alice.New(
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
	)

	dbJobsEngine := launchBGJobEngine(environment, session)
	defer dbJobsEngine.Stop() // does it leak mongo session?

	log.Printf("All startup routines completed successfully, app is listening on %s port\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), defaultMiddleware.Then(router)))
}

func launchBGJobEngine(env *utils.Environment, session *mgo.Session) *cron.Cron {
	log.Println("Setting up and launching the background jobs engine")
	bgJobEngine := cron.New()

	// Runs 1/2 hourly at the beginning of each hour and at 30th minute of each hour
	// ---------------- s  m   h d m
	bgJobEngine.AddJob("0 0,30 * * *", jobs.NewStopTimersAtMidnight(env, session.Clone()))
	log.Println("--- Scheduled StopTimersAtMidnight job")

	// Runs once an hour at 15 minutes
	// ---------------- s  m   h d m
	bgJobEngine.AddJob("0 25 * * *", jobs.NewClearPasses(env, session.Clone()))
	log.Println("--- Scheduled ClearPasses job")

	bgJobEngine.Start()
	return bgJobEngine
}

func getEnvironmentName() string {
	env := os.Getenv("SLACK_TIME_ENV")
	if env == "" {
		env = utils.DevelopmentEnv
	}
	return env
}
