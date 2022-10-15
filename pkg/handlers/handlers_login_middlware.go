package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

func (handlers *Handlers) LoginMiddlware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			cookie, err := r.Cookie("Access-token")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"result": "error", "message": "invalid token"})
				return
			}
			user := &models.User{
				Jwt: cookie.Value,
			}
			if err := handlers.service.ValidateToken(user); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"result": "error", "message": "invalid token"})
				return
			}
			if user.CsrfToken != r.Header.Get("X-CSRF-Token") {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"result": "error", "message": "invalid token"})
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(context.TODO(), KeyForToken("user_id"),user.UserId.String())))
		})
	}
}
