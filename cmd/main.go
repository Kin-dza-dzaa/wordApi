package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/apierror"
	"github.com/Kin-dza-dzaa/wordApi/pkg/handlers"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	service "github.com/Kin-dza-dzaa/wordApi/pkg/servise"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

func main() {
	myLogger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	config, err := config.ReadConfig()
	if err != nil {
		myLogger.Fatal().Msg(err.Error())
	}
	pool, err := pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		myLogger.Fatal().Msg(err.Error())
	}
	defer pool.Close()
	myRepository := repositories.NewRepository(pool)
	myService := service.NewService(myRepository, config)
	myApiError := apierror.NewApiError(&myLogger)
	myHandlers := handlers.NewHandlers(myService, config, myApiError)
	srv := &http.Server{
		Handler:      myHandlers.Cors.Handler(myHandlers.Router),
		Addr:         config.Adress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		myLogger.Info().Msg(fmt.Sprintf("Staring server wordapi at %v", config.Adress))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			myLogger.Fatal().Msg(err.Error())
		}
	}()
	<-myHandlers.StopHTTPServerChan
	if err := srv.Shutdown(context.TODO()); err != nil {
		myLogger.Fatal().Msg(err.Error())
	}
}
