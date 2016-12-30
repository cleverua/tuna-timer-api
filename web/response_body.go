package web

import (
	"github.com/cleverua/tuna-timer-api/models"
)

// Common response body for frontend application
type ResponseBody struct {
	AppInfo        map[string]string `json:"appInfo"`
	ResponseStatus *ResponseStatus	 `json:"response_status"`
	ResponseData   map[string]string `json:"data"`
}

type ResponseStatus struct {
	Status		 string `json:"status"`
	DeveloperMessage string `json:"developerMessage"`
	UserMessage	 string `json:"user_message"`
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

// Response with array of projects data
type ProjectsResponseBody struct {
	ResponseBody
	ResponseData []*models.Project `json:"data"`
}
