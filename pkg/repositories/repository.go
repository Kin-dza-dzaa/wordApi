package repositories

import (
	"context"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddWords(ctx context.Context, words models.WordsToAdd, badwords map[string]string, userId string) error
	GetWords(ctx context.Context, words models.Words, userId string) error
	UpdateWord(ctx context.Context, words models.WordToUpdate, userId string) error
	DeleteWords(ctx context.Context, words models.WordsToDelete, userId string) error
	UpdateState(ctx context.Context, words models.StatesToUpdate, userId string) error
	IfWordInDb(ctx context.Context, word string, result *bool) error
	IfUserHasWord(ctx context.Context, word string, collectionName string, result *bool, userId string) error
}

func NewRepository(pool *pgxpool.Pool) Repository{
	return NewRepositoryWord(pool)
}