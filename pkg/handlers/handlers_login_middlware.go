package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handlers) LoginMiddlware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, v := range loginRoutes {
				if r.URL.Path == v {
					token := strings.Split(r.Header.Get("Authorization"), " ")
					if len(token) != 2 {
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
						return
					}
					user_id, err := h.service.ValidateToken(token[1])
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
						return
					}
					next.ServeHTTP(w, r.WithContext(context.WithValue(context.TODO(), KEY, user_id)))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
