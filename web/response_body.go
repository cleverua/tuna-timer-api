package web

import (
	"github.com/cleverua/tuna-timer-api/models"
)

// Common response body for frontend application
type ResponseBody struct {
	AppInfo        map[string]string `json:"app_info"`
	ResponseStatus *ResponseStatus	 `json:"response_status"`
	ResponseData   map[string]string `json:"data"`
}

func NewResponseBody(info map[string]string) *ResponseBody{
	return &ResponseBody{
		ResponseStatus: &ResponseStatus{ Status: statusOK },
		AppInfo: info,
	}
}

type ResponseStatus struct {
	Status		 string `json:"status"`
	DeveloperMessage string `json:"developer_message"`
	UserMessage	 string `json:"user_message"`
}

// Response with jwt token
type JWTResponse struct {
	*ResponseBody
	ResponseData JwtToken `json:"data"`
}

func NewJWTResponse(info map[string]string) *JWTResponse {
	return &JWTResponse{
		ResponseData: JwtToken{},
		ResponseBody: NewResponseBody(info),
	}
}

// Response with task data
type TimerResponse struct {
	*ResponseBody
	ResponseData models.Timer	`json:"data"`
}

func NewTimerResponse(info map[string]string) *TimerResponse {
	return &TimerResponse{
		ResponseBody: NewResponseBody(info),
	}
}

// Response with array of tasks data
type TimersResponse struct {
	*ResponseBody
	ResponseData []*models.Timer `json:"data"`
}

func NewTimersResponse(info map[string]string) *TimersResponse {
	return &TimersResponse{
		ResponseBody: NewResponseBody(info),
	}
}

// Response with array of projects data
type ProjectsResponse struct {
	*ResponseBody
	ResponseData []*models.Project `json:"data"`
}

func NewProjectsResponse(info map[string]string) *ProjectsResponse {
	return &ProjectsResponse{
		ResponseBody: NewResponseBody(info),
	}
}

// Response with array of user month statistics data
type UserStatisticsResponse struct {
	*ResponseBody
	ResponseData []*models.UserStatisticsAggregation `json:"data"`
}

func NewUserStatisticsResponse(info map[string]string) *UserStatisticsResponse {
	return &UserStatisticsResponse{
		ResponseBody: NewResponseBody(info),
	}
}
