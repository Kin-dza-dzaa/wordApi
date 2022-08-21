package handlers

import (
	"github.com/gorilla/mux"
	"github.com/Kin-dza-dzaa/wordApi/pkg/servise"
)

type Handlers struct {
	service service.Service
	Router *mux.Router
}

func (h *Handlers) InitilizeHandlers() {
	h.Router = mux.NewRouter().Host("localhost").Subrouter()
	h.Router.Handle("/sign-up", h.SignUpUser()).Methods("POST")
}

func NewHandlers(servce service.Service) *Handlers{
	return &Handlers{
		service: servce,
	}
}