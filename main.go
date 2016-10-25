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

	public := alice.New(web.LoggingMiddleware, web.RecoveryMiddleware)
	secure := alice.New(web.LoggingMiddleware, web.RecoveryMiddleware, web.JWTMiddleware)


	router := mux.NewRouter().StrictSlash(true)


	router.Handle("/api/v1/health", public.ThenFunc(handlers.Health)).Methods("GET")

	// Slack will sometimes call the API method using a GET request
	// to check SSL certificate - so we reply with a status handler here
	router.Handle("/api/v1/timer", public.ThenFunc(handlers.Timer)).Methods("POST", "GET")

	// Slack  OAuth2 stuff
	router.Handle("/api/v1/slack/oauth2redirect", public.ThenFunc(handlers.SlackOauth2Redirect)).Methods("GET")

	// Static assets
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	// ===== Routes for frontend application
	// Activates the pass and returns back a JWT token, it essentially logs the user in
	router.Handle("/api/v1/frontend/auth/{token}/activate", public.ThenFunc(handlers.NotImplemented)).Methods("POST")
	// reads JWT from header and returns a 200 if it is okay and not expired
	router.Handle("/api/v1/frontend/auth/validate", secure.ThenFunc(handlers.ValidateAuthToken)).Methods("GET")


	// Temporary stuff, remove eventually
	router.Handle("/api/v1/temporary/clear_data", public.ThenFunc(handlers.ClearAllData)).Methods("GET")
	router.Handle("/api/v1/temporary/send_message", public.ThenFunc(handlers.SendSampleMessageFromBot)).Methods("GET")
	router.Handle("/api/v1/temporary/new_jwt_token", public.ThenFunc(handlers.NewJWTToken)).Methods("GET")

	dbJobsEngine := launchBGJobEngine(environment, session)
	defer dbJobsEngine.Stop() // does it leak mongo session?

	log.Printf("All startup routines completed successfully, app is listening on %s port\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
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
