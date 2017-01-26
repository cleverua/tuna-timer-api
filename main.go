package main

import (
	"net/http"
	"os"

	"log"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/cleverua/tuna-timer-api/utils"
	"github.com/cleverua/tuna-timer-api/web"
	"time"

	"fmt"
	"github.com/robfig/cron"
	"github.com/cleverua/tuna-timer-api/jobs"
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
	fh := web.NewFrontendHandlers(environment, session)
	secureCTX := web.SecureContext{
		Origin:  environment.Config.UString("origin.url"),
		Session: session,
		Env: 	 environment,
	}

	public := alice.New(web.LoggingMiddleware, web.RecoveryMiddleware, secureCTX.CorsMiddleware)
	secure := alice.New(
		web.LoggingMiddleware,
		web.RecoveryMiddleware,
		secureCTX.CorsMiddleware,
		web.JWTMiddleware,
		secureCTX.CurrentUserMiddleware)

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
	router.Handle("/api/v1/frontend/session", public.ThenFunc(fh.Authenticate)).Methods("POST", "OPTIONS")
	// Routes for user data CRUD
	router.Handle("/api/v1/frontend/timers", secure.ThenFunc(fh.TimersData)).Methods("GET", "OPTIONS")
	router.Handle("/api/v1/frontend/timers", secure.ThenFunc(fh.CreateTimer)).Methods("POST", "OPTIONS")
	router.Handle("/api/v1/frontend/timers/{id}", secure.ThenFunc(fh.UpdateTimer)).Methods("PUT", "OPTIONS")
	router.Handle("/api/v1/frontend/timers/{id}", secure.ThenFunc(fh.DeleteTimer)).Methods("DELETE", "OPTIONS")
	router.Handle("/api/v1/frontend/projects", secure.ThenFunc(fh.ProjectsData)).Methods("GET", "OPTIONS")
	router.Handle("/api/v1/frontend/month_statistics", secure.ThenFunc(fh.MonthStatistic)).Methods("GET", "OPTIONS")

	// Temporary stuff, remove eventually
	router.Handle("/api/v1/frontend/auth/validate", secure.ThenFunc(handlers.ValidateAuthToken)).Methods("GET", "OPTIONS")
	router.Handle("/api/v1/temporary/clear_data", public.ThenFunc(handlers.ClearAllData)).Methods("GET")
	router.Handle("/api/v1/temporary/send_message", public.ThenFunc(handlers.SendSampleMessageFromBot)).Methods("GET")

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
