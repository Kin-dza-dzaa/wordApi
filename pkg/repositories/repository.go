package repositories

import (
	"github.com/jackc/pgx/v4"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

type RepositoryUser interface {
	SignUpUser(user *models.User) error
	GetUser(user *models.User) (*models.User, error)
}

type RepositoryWord interface {
	AddWords(words models.Words, userId string)
	GetWords(userId string) (*models.Words, error)
	UpdateWord(words models.Words, userId string) error
	DeleteWords(words models.Words, userId string)
}

type Repository struct {
	RepositoryUser
	RepositoryWord
}

func NewRepository(conn *pgx.Conn) *Repository{
	return &Repository{
		RepositoryUser: &repositoryUser{conn: conn}, 
		RepositoryWord: &repositoryWord{conn: conn},
	}
}
