package web

import (
	"net/http"
	"github.com/tuna-timer/tuna-timer-api/data"
	"encoding/json"
)

func (h *Handlers) UserAuthentication(w http.ResponseWriter, r *http.Request) {
	pid := r.PostFormValue("pid")

	session := h.mongoSession.Clone()
	defer session.Close()

	pass_service := data.NewPassService(session)
	pass, err := pass_service.FindPassByToken(pid)
	result := make(map[string]string)

	if err == nil && pass == nil {
		w.WriteHeader(http.StatusBadRequest)
		result["userMessage"] = "please login from slack application"
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result["developerMessage"] = err.Error()
	} else {
		jwt_token, jwt_err := NewToken(pass, session)
		if jwt_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result["developerMessage"] = jwt_err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
			result["jwt"] = string(jwt_token)
		}
	}
	json.NewEncoder(w).Encode(result)
}
