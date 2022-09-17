package repositories

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4"
)

type RepositoryUser interface {
	SignUpUser(user *models.User) error
	GetUser(user *models.User) (*models.User, error)
}

type RepositoryWord interface {
	AddWords(words models.WordsAdd, userId string) []string
	GetWords(userId string) (*models.WordsGet, error)
	UpdateWord(words models.WordsUpdate, userId string) error
	DeleteWords(words models.WordsDelete, userId string)
}

type Repository struct {
	RepositoryUser
	RepositoryWord
}

func NewRepository(conn *pgx.Conn, config *config.Config) *Repository{
	return &Repository{
		RepositoryUser: &repositoryUser{conn: conn}, 
		RepositoryWord: &repositoryWord{conn: conn, config: config},
	}
}
