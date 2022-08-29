package main

import (
	"context"
	"log"
	"net/http"
	"time"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/logger"
	"github.com/Kin-dza-dzaa/wordApi/pkg/handlers"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/wordApi/pkg/servise"
)

func main() {
	handlers.StopHTTPServerChan = make(chan bool)
	w, err := logger.GetWriter()
	if err != nil {
		log.Fatal(err)
	}
	//myLogger := logger.Getlogger(w)
	defer func() {
		if err := w.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := repositories.Connect(config.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	repository := repositories.NewRepository(conn, config)
	service := service.NewService(repository, config)
	handler := handlers.NewHandlers(service)
	handler.InitilizeHandlers()
	srv := &http.Server{
		Handler: handler.Router,
		Addr:    "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	<-handlers.StopHTTPServerChan
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
