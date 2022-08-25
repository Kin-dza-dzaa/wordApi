package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

func (h *Handlers) SignUpUserHandler() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		var userInput models.User
		w.Header().Set("Content-type", "application/json")
		if json.NewDecoder(r.Body).Decode(&userInput) != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "Wrong input, expected email, user_name and password"})
			return
		}
		if err := h.service.SignUpUser(&userInput); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return 
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "Successfully signed up"})
	})
}

func (h *Handlers) SignInUserHandler() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		var userInput models.User
		w.Header().Set("Content-type", "application/json")
		if json.NewDecoder(r.Body).Decode(&userInput) != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "Wrong input, expected email, user_name and password"})
			return
		}
		token, err := h.service.SignInUser(&userInput)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		w.Header().Set("Authorization", "Bearer " + token)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "token was sent"})
	})
}
