package service

import (
	"context"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
)

type Service interface {
	AddWords(ctx context.Context, words models.WordsToAdd, userId string) (map[string]string, error)
	GetWords(ctx context.Context, words models.Words, userId string) error
	UpdateWord(ctx context.Context, words models.WordToUpdate, userId string) error
	DeleteWords(ctx context.Context, words models.WordsToDelete, userId string) error
	UpdateState(ctx context.Context, words models.StatesToUpdate, userId string) error
	ValidateToken(user *models.User) error
}

func NewService(repository repositories.Repository, config *config.Config) Service{
	return NewServiceWords(repository, config)
}