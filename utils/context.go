package utils

import (
	"context"
	"gopkg.in/mgo.v2"
)

const (
	ContextKeyMongoSession = "mongoSession"
	ContextKeySelfBaseURL  = "selfBaseUrl"
	ContextKeyTheme        = "theme"
)

func GetMongoSessionFromContext(ctx context.Context) *mgo.Session {
	return ctx.Value(ContextKeyMongoSession).(*mgo.Session)
}

func PutMongoSessionInContext(parentContext context.Context, session *mgo.Session) context.Context {
	return context.WithValue(parentContext, ContextKeyMongoSession, session)
}

func PutSelfBaseURLInContext(parentContext context.Context, url string) context.Context {
	return context.WithValue(parentContext, ContextKeySelfBaseURL, url)
}

func GetSelfBaseURLFromContext(ctx context.Context) string {
	return ctx.Value(ContextKeySelfBaseURL).(string)
}

func PutThemeInContext(parentContext context.Context, theme interface{}) context.Context {
	return context.WithValue(parentContext, ContextKeyTheme, theme)
}

func GetThemeFromContext(ctx context.Context) interface{} {
	return ctx.Value(ContextKeyTheme)
}
