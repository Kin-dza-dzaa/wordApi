package handlers

import (
	"strings"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/apierror"
	service "github.com/Kin-dza-dzaa/wordApi/pkg/servise"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type KeyForToken string

type Handlers struct {
	StopHTTPServerChan chan bool
	service            service.Service
	Router             *mux.Router
	Cors               *cors.Cors
	config             *config.Config
	ApiError           *apierror.ApiError
}

func NewHandlers(service service.Service, config *config.Config, ApiError *apierror.ApiError) *Handlers {
	handlers := new(Handlers)
	handlers.config = config
	handlers.ApiError = ApiError
	handlers.service = service
	handlers.StopHTTPServerChan = make(chan bool)
	handlers.Cors = cors.New(cors.Options{
		AllowedOrigins:   strings.Split(handlers.config.AllowedOrigns, ","),
		AllowCredentials: config.AllowCredentials,
		AllowedHeaders:   []string{"User-Agent", "Content-type", "X-Csrf-Token"},
		MaxAge:           5,
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})
	handlers.Router = mux.NewRouter()
	words := handlers.Router.PathPrefix("/words").Subrouter()
	words.Handle("", handlers.ApiError.Middleware(handlers.GetWordsHandler())).Methods("GET")
	words.Handle("", handlers.ApiError.Middleware(handlers.DeleteWordHandler())).Methods("DELETE")
	words.Handle("", handlers.ApiError.Middleware(handlers.UpdateWordHandler())).Methods("PUT")
	words.Handle("", handlers.ApiError.Middleware(handlers.AddWordsHandler())).Methods("POST")
	words.Handle("/state", handlers.ApiError.Middleware(handlers.UpdateStateHandler())).Methods("PUT")

	handlers.Router.Use(handlers.LoginMiddleware())
	return handlers
}
