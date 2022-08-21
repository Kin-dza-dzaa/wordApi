package models

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

var Conn *pgx.Conn

func Connect(url string, logger *zerolog.Logger) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		logger.Error().Msg(err.Error())
		return
	}
	Conn = conn
}
