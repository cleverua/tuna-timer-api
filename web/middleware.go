package web

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
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

type Cors struct {
	Origin string
}

func (c Cors) CorsMiddleware(h http.Handler) http.Handler {
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
