package web

import (
	"github.com/cleverua/tuna-timer-api/models"
)

// Common response body for frontend application
type ResponseBody struct {
	AppInfo        map[string]string `json:"appInfo"`
	ResponseErrors map[string]string `json:"errors"`
	ResponseData   map[string]string `json:"data"`
}

// Response with jwt token
type JwtResponseBody struct {
	ResponseBody
	ResponseData JwtToken `json:"data"`
}

// Response with task data
type TaskResponseBody struct {
	ResponseBody
	ResponseData models.TaskAggregation `json:"data"`
	TaskErrors   map[string]string      `json:"errors"`
}

// Response with array of tasks data
type TasksResponseBody struct {
	ResponseBody
	ResponseData []*models.Timer `json:"data"`
}
