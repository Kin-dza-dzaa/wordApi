package repository

import (
	"github.com/jackc/pgx/v4"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

type Authorization interface {
	SignUpUser(user *models.User) error
	GetUser(email string) (*models.User, error)
}

type Repository interface {
	Authorization
}

func NewRepository(conn *pgx.Conn) Repository{
	return &repository{conn: conn}
}
