package service

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repository"
	"github.com/google/uuid"
)

type Authorization interface {
	SignInUser(user *models.User) (uuid.UUID, error)
	GenerateToken(user *models.User) (string, error)
	ValidateToken(token string) (string, error)
	SignUpUser(user *models.User) error
}

type Service interface {
	Authorization
}

func NewService(repository repository.Repository, config *config.Config) Service{
	return &service{repository: repository, config: config}
	
}
