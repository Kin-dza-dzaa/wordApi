package repositories

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4"
)

type RepositoryUser interface {
	SignUpUser(user *models.User) error
	GetUser(user *models.User) (*models.User, error)
}

type RepositoryWord interface {
	AddWords(words models.Words, userId string) []string
	GetWords(userId string) (*[]external.Translation, error)
	UpdateWord(words models.Words, userId string) error
	DeleteWords(words models.Words, userId string)
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
