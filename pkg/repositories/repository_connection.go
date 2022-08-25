package repositories

import (
	"github.com/jackc/pgx/v4"
	"context"
	
)

func Connect(url string) (*pgx.Conn, error){
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
