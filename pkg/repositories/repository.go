package repositories

import (
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Repository interface {
	AddWords(words *models.WordsToAdd, userId string) []string
	GetWords(userId string) (*models.WordsGet, error)
	UpdateWord(words *models.WordToUpdate, userId string) error
	DeleteWords(words *models.WordsToDelete, userId string)
	UpdateState(words *models.StatesToUpdate, userId string)
}

func NewRepository(pool *pgxpool.Pool, logger *zerolog.Logger, config *config.Config) Repository{
	return NewRepositoryWord(pool, logger, config)
}
