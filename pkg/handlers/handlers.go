package handlers

import (
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/pkg/servise"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type keyForToken string

var (
	StopHTTPServerChan chan bool
	loginRoutes []string = []string{"/words", "/user/log-out", "/user/check"}
)

const (
	KEY keyForToken = "user_id"
)

type Handlers struct {
	service *service.Service
	Router *mux.Router
	Cors *cors.Cors
}

func (h *Handlers) ShutDown() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		StopHTTPServerChan <- true
	})
}

func (h *Handlers) InitilizeHandlers() {
	h.Cors = cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedHeaders: []string{"User-Agent", "Content-type"},
		MaxAge: 5,
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
	})
	h.Router = mux.NewRouter()
	h.Router.Handle("/user", h.SignUpUserHandler()).Methods("POST")
	h.Router.Handle("/user/token", h.SignInUserHandler()).Methods("POST")
	h.Router.Handle("/user/log-out", h.LogOutUser()).Methods("GET")
	h.Router.Handle("/user/check", h.CheckUser()).Methods("GET")
	h.Router.Handle("/words", h.GetWordsHandler()).Methods("GET")
	h.Router.Handle("/words", h.DeleteWordHandler()).Methods("DELETE")
	h.Router.Handle("/words", h.UpdateWordHandler()).Methods("PUT")
	h.Router.Handle("/words", h.AddWordsHandler()).Methods("POST")
	h.Router.Handle("/", h.ShutDown()).Methods("GET")
	h.Router.Use(h.LoginMiddlware())	
}

func NewHandlers(service *service.Service) *Handlers{
	return &Handlers{
		service: service,
	}
}
