package handlers

import (
	service "github.com/Kin-dza-dzaa/wordApi/pkg/servise"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type KeyForToken string

type Handlers struct {
	StopHTTPServerChan chan bool
	service service.Service
	Router  *mux.Router
	Cors    *cors.Cors
}

func (h *Handlers) InitilizeHandlers() {
	h.StopHTTPServerChan = make(chan bool)
	h.Cors = cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"User-Agent", "Content-type"},
		MaxAge:           5,
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})
	h.Router = mux.NewRouter()
		words := h.Router.PathPrefix("/words").Subrouter()
		words.Handle("", h.GetWordsHandler()).Methods("GET")
		words.Handle("", h.DeleteWordHandler()).Methods("DELETE")
		words.Handle("", h.UpdateWordHandler()).Methods("PUT")
		words.Handle("", h.AddWordsHandler()).Methods("POST")
		words.Handle("/state", h.UpdateStateHandler()).Methods("PUT")
	
	h.Router.Use(h.LoginMiddlware())
}

func NewHandlers(service service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}
