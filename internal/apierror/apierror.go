package apierror

import (
	"errors"
	"net/http"

	"github.com/jackc/puddle"
	"github.com/rs/zerolog"
)

type HttpErrHandler func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Logger *zerolog.Logger
}

func (apiError *ApiError) Middleware(next HttpErrHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			switch err := err.(type) {

			case *Response:
				w.WriteHeader(err.StatusCode)
				w.Write(err.Marshal())

			default:
				apiError.Logger.Error().Msg(err.Error())
				if errors.Is(err, puddle.ErrClosedPool) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(NewResponse("error", "internal server error", http.StatusInternalServerError).Marshal())
					return
				}

				if errors.Is(err, puddle.ErrNotAvailable) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(NewResponse("too many requests", "error", http.StatusTooManyRequests).Marshal())
					return
				}

				w.Write(NewResponse("error", "unexpected error", http.StatusInternalServerError).Marshal())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})
}

func NewApiError(logger *zerolog.Logger) *ApiError {
	return &ApiError{
		Logger: logger,
	}
}
