package web

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pavlo/slack-time/utils"
)

// Handlers is a collection of net/http handlers to serve the API
type Handlers struct {
	env    *utils.Environment
	status map[string]string
}

// NewHandlers constructs a Handlers collection
func NewHandlers(env *utils.Environment) *Handlers {
	return &Handlers{env: env, status: map[string]string{
		"env":     env.Name,
		"version": env.AppVersion,
	}}
}

// Timer handles Slack /timer command
func (h *Handlers) Timer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

// Health handles a call for app health request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.env.CreatedAt)
	h.status["uptime"] = uptime.String() //is it good or not if I modify the map here?
	json.NewEncoder(w).Encode(h.status)
}
