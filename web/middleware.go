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
