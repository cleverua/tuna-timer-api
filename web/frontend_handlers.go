package web

import (
	"net/http"
	"github.com/cleverua/tuna-timer-api/data"
	"encoding/json"
)

func (h *Handlers) UserAuthentication(w http.ResponseWriter, r *http.Request) {
	response := JwtResponseBody{
		ResponseData: JwtToken{},
		ResponseBody: ResponseBody{
			ResponseErrors: map[string]string{},
			AppInfo: h.status,
		},
	}
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
		jwt_token, jwt_err := NewToken(pass.TeamUserID, session)
		if jwt_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseErrors["developerMessage"] = jwt_err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
			response.ResponseData.Token = jwt_token
		}
	}
	json.NewEncoder(w).Encode(response)
}
