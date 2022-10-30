package handlers

import (
	"context"
	"net/http"

	"github.com/Kin-dza-dzaa/wordApi/internal/apierror"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

func (handlers *Handlers) LoginMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			cookie, err := r.Cookie("Access-token")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(apierror.NewResponse("error", "cookie not present", http.StatusBadRequest).Marshal())
				return
			}
			user := &models.User{
				Jwt: cookie.Value,
			}
			if err := handlers.service.ValidateToken(user); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(apierror.NewResponse("error", "invalid token", http.StatusBadRequest).Marshal())
				return
			}
			if user.CsrfToken != r.Header.Get("X-CSRF-Token") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(apierror.NewResponse("error", "XSRF failed", http.StatusBadRequest).Marshal())
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(context.TODO(), KeyForToken("user_id"), user.UserId.String())))
		})
	}
}
