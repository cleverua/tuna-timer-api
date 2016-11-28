package utils

import (
	"crypto/tls"
	"github.com/cleverua/tuna-timer-api/models"
	"net/http"
	"testing"
	"gopkg.in/tylerb/is.v1"
)

func TestNormalizeSlackCustomCommand(t *testing.T) {

	s := is.New(t)

	cmd := NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "start Add MongoDB service to docker-compose.yml",
	})

	s.Equal(cmd.Text, "Add MongoDB service to docker-compose.yml")
	s.Equal(cmd.SubCommand, "start")

	cmd = NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "      start         Add MongoDB service to docker-compose.yml",
	})

	s.Equal(cmd.Text, "Add MongoDB service to docker-compose.yml")
	s.Equal(cmd.SubCommand, "start")

	cmd = NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "status",
	})

	s.Equal(cmd.Text, "")
	s.Equal(cmd.SubCommand, "status")

	cmd = NormalizeSlackCustomCommand(models.SlackCustomCommand{
		Text: "start ",
	})
	s.Equal(cmd.Text, "")
	s.Equal(cmd.SubCommand, "start")
}

func TestGetSelfURLFromRequest(t *testing.T) {
	s := is.New(t)

	r := &http.Request{
		Host: "subdomain.domain.com",
	}

	s.Equal(GetSelfURLFromRequest(r), "http://subdomain.domain.com")

	r = &http.Request{
		Host: "subdomain.domain.com",
		TLS:  &tls.ConnectionState{},
	}

	s.Equal(GetSelfURLFromRequest(r), "https://subdomain.domain.com")
}

