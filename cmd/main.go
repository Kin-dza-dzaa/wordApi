package main

import (
	"context"
	"log"
	"net/http"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/logger"
	"github.com/Kin-dza-dzaa/wordApi/pkg/handlers"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repository"
	service "github.com/Kin-dza-dzaa/wordApi/pkg/servise"
)

func main() {
	w, err := logger.GetWriter()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	//myLogger := logger.Getlogger(w)
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := repository.Connect(config.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	repository := repository.NewRepository(conn)
	service := service.NewService(repository, config)
	handler := handlers.NewHandlers(service)
	handler.InitilizeHandlers()
	http.ListenAndServe(":8080", handler.Router)
	defer conn.Close(context.Background())
}
