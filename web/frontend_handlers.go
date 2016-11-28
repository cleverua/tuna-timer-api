package web

import (
	"net/http"
	"github.com/cleverua/tuna-timer-api/data"
	"encoding/json"
)

type ResponseBody struct {
	AppInfo map[string]string `json:"appInfo"`
	ResponseErrors map[string]string `json:"errors"`
	ResponseData interface{} `json:"data"`
}

func NewResponseBody(h *Handlers) *ResponseBody {
	return &ResponseBody{
		ResponseErrors: map[string]string{},
		AppInfo: h.status,
	}
}

func (h *Handlers) UserAuthentication(w http.ResponseWriter, r *http.Request) {
	response := NewResponseBody(h)
	pid := r.PostFormValue("pid") // TODO: sanitize pid
	session := h.mongoSession.Clone()
	defer session.Close()

	pass_service := data.NewPassService(session)
	pass, err := pass_service.FindPassByToken(pid)

	if err == nil && pass == nil {
		w.WriteHeader(http.StatusBadRequest)
		response.ResponseErrors["userMessage"] = "please login from slack application"
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ResponseErrors["developerMessage"] = err.Error()
	} else {
		jwt_token, jwt_err := NewToken(pass, session)
		if jwt_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseErrors["developerMessage"] = jwt_err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
			response.ResponseData = map[string]string{"jwt": jwt_token}
		}
	}
	json.NewEncoder(w).Encode(response)
}
