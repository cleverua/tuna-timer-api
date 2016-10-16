package utils

import (
	"testing"
	"gopkg.in/tylerb/is.v1"
)


func TestNewEnvironment(t *testing.T) {
	s := is.New(t)
	env := NewEnvironment(TestEnv, "1")

	s.Equal(env.AppVersion, "1")
	s.NotNil(env.CreatedAt) //todo - check type rather?
}
