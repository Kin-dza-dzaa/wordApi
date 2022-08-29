package service

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
)

type ServiceUser interface {
	SignInUser(user *models.User) (string, error)
	ValidateToken(token string) (string, error)
	SignUpUser(user *models.User) error
}

type ServiceWord interface {
	AddWords(words models.Words, userId string) []string
	GetWords(userId string) (*[]external.Translation, error)
	UpdateWord(words models.Words, userId string) error
	DeleteWords(words models.Words, userId string)
}

type Service struct {
	ServiceUser
	ServiceWord
}

func NewService(repository *repositories.Repository, config *config.Config) *Service{
	return &Service{
		ServiceUser: &serviceUser{repository: repository, config: config},
		ServiceWord: &serviceWord{repository: repository},
	}
}
