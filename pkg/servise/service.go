package service

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/go-playground/validator/v10"
)

type ServiceUser interface {
	SignInUser(user *models.User) (string, error)
	ValidateToken(token string) (string, error)
	SignUpUser(user *models.User) error
}

type ServiceWord interface {
	AddWords(words models.WordsAdd, userId string) []string
	GetWords(userId string) (*models.WordsGet, error)
	UpdateWord(words models.WordToUpdate, userId string) error
	DeleteWords(words models.WordsDelete, userId string)
}

type Service struct {
	ServiceUser
	ServiceWord
}

func NewService(repository *repositories.Repository, config *config.Config, validator *validator.Validate) *Service{
	return &Service{
		ServiceUser: &serviceUser{repository: repository, config: config, validator: validator},
		ServiceWord: &serviceWord{repository: repository},
	}
}
