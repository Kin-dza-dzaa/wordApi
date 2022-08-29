package handlers

import (
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/pkg/servise"
	"github.com/gorilla/mux"
)

type keyForToken string

var (
	key keyForToken = "user_id"
	StopHTTPServerChan chan bool
	loginRoutes []string = []string{"/words/delete", "/words/update", "/words/add", "/words"}
)

type Handlers struct {
	service *service.Service
	Router *mux.Router
}

func (h *Handlers) ShutDown() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		StopHTTPServerChan <- true
	})
}

func (h *Handlers) InitilizeHandlers() {
	h.Router = mux.NewRouter().Host("localhost").Subrouter()
	h.Router.Handle("/sign-up", h.SignUpUserHandler()).Methods("POST")
	h.Router.Handle("/sign-in", h.SignInUserHandler()).Methods("POST")
	h.Router.Handle("/words", h.GetWordsHandler()).Methods("GET")
	h.Router.Handle("/words/delete", h.DeleteWordHandler()).Methods("DELETE")
	h.Router.Handle("/words/update", h.UpdateWordHandler()).Methods("PUT")
	h.Router.Handle("/words/add", h.AddWordsHandler()).Methods("POST")
	h.Router.Handle("/shut-down", h.ShutDown()).Methods("GET")
	h.Router.Use(h.LoginMiddlware())	
}

func NewHandlers(service *service.Service) *Handlers{
	return &Handlers{
		service: service,
	}
}
