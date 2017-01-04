package web

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"github.com/cleverua/tuna-timer-api/data"
	"github.com/gorilla/context"
)

func LoggingMiddleware(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

func RecoveryMiddleware(h http.Handler) http.Handler {
	return handlers.RecoveryHandler()(h)
}

func JWTMiddleware(h http.Handler) http.Handler {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("TODO: Extract me in config/env"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	}).Handler(h)
}

type SecureContext struct {
	Origin	string
	Session *mgo.Session
}

func (c *SecureContext) CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == c.Origin {
			w.Header().Set("Access-Control-Allow-Origin", c.Origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Origin, Authorization")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (c *SecureContext) CurrentUserMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData := context.Get(r, "user").(*jwt.Token).Claims.(jwt.MapClaims)

		userService := data.NewUserService(c.Session)
		user, err := userService.FindByID(userData["user_id"].(string))
		if err != nil {
			response := ResponseBody{}
			response.ResponseStatus.Status = statusBadRequest
			response.ResponseStatus.UserMessage = userLoginMessage
			response.ResponseStatus.DeveloperMessage = err.Error()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		context.Set(r, "user", user)
		h.ServeHTTP(w, r)
	})
}
