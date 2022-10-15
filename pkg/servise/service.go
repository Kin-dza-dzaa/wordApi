package service

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/rs/zerolog"
)

type Service interface {
	AddWords(words *models.WordsToAdd, userId string) []string
	GetWords(userId string) (*models.WordsGet, error)
	UpdateWord(words *models.WordToUpdate, userId string) error
	DeleteWords(words *models.WordsToDelete, userId string)
	UpdateState(words *models.StatesToUpdate, userId string)
	ValidateToken(user *models.User) (error)
}

func NewService(repository repositories.Repository, config *config.Config, logger *zerolog.Logger) Service{
	return NewServiceWords(repository, config, logger)
}
