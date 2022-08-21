package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

func (h *Handlers) SignUpUser() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		var userInput models.User
		w.Header().Set("Content-type", "Application/json")
		if json.NewDecoder(r.Body).Decode(&userInput) != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "Wrong input, expected email, user_name and password"})
			return
		}
		if h.service.SignUpUser(&userInput) != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "User already exists"})
			return 
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "Successfully signed up"})
	})
}