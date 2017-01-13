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
type JWTResponseBody struct {
	ResponseBody
	ResponseData JwtToken `json:"data"`
}

func NewJWTResponseBody(info map[string]string) *JWTResponseBody{
	return &JWTResponseBody{
		ResponseData: JwtToken{},
		ResponseBody: ResponseBody{
			ResponseStatus: &ResponseStatus{ Status: statusOK },
			AppInfo: info,
		},
	}
}

// Response with task data
type TimerResponseBody struct {
	ResponseBody
	ResponseData models.Timer	`json:"data"`
	TaskErrors   map[string]string  `json:"errors"`
}

func NewTimerResponseBody(info map[string]string) *TimerResponseBody {
	return &TimerResponseBody{
		ResponseBody: ResponseBody{
			ResponseStatus: &ResponseStatus{ Status: statusOK },
			AppInfo: info,
		},
		TaskErrors: map[string]string{},
	}
}

// Response with array of tasks data
type TimersResponseBody struct {
	ResponseBody
	ResponseData []*models.Timer `json:"data"`
}

func NewTimersResponseBody(info map[string]string) *TimersResponseBody {
	return &TimersResponseBody{
		ResponseBody: ResponseBody{
			ResponseStatus: &ResponseStatus{ Status: statusOK },
			AppInfo: info,
		},
	}
}

// Response with array of projects data
type ProjectsResponseBody struct {
	ResponseBody
	ResponseData []*models.Project `json:"data"`
}

func NewProjectsResponseBody (info map[string]string) *ProjectsResponseBody {
	return &ProjectsResponseBody{
		ResponseBody: ResponseBody{
			ResponseStatus: &ResponseStatus{ Status: statusOK },
			AppInfo: info,
		},
	}
}
