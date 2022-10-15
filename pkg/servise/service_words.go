package service

import (
	"errors"
	"time"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SercviceWords struct {
	repository repositories.Repository
	config *config.Config
	logger *zerolog.Logger
}

func (service *SercviceWords) AddWords(words *models.WordsToAdd, userId string) []string {
	for i := 0; i < len(words.Words); i++ {
		words.Words[i].TimeOfLastRepeating = time.Now().UTC()
	}
	return service.repository.AddWords(words, userId)
}

func (service *SercviceWords) GetWords(userId string) (*models.WordsGet, error) {
	return service.repository.GetWords(userId)
}

func (service *SercviceWords) UpdateWord(words *models.WordToUpdate, userId string) error {
	words.TimeOfLastRepeating = time.Now().UTC()
	return service.repository.UpdateWord(words, userId)
}

func (service *SercviceWords) DeleteWords(words *models.WordsToDelete, userId string) {
	service.repository.DeleteWords(words, userId)
}

func (service *SercviceWords) UpdateState(words *models.StatesToUpdate, userId string) {
	service.repository.UpdateState(words, userId)
}

func (service *SercviceWords) ValidateToken(user *models.User) (error) {
	token, err := jwt.ParseWithClaims(user.Jwt, &models.MyJwtClaims{}, func(t *jwt.Token) (interface{}, error) {return []byte(service.config.JWTString), nil})
	if err != nil {
		service.logger.Warn().Msg(err.Error())
		return errors.New("invalid token")
	}
	claims, ok := token.Claims.(models.MyJwtClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token")
	}
	user.CsrfToken = claims.XCSRFToken
	Uuid, err := uuid.Parse(claims.UserId)
	if err != nil {
		return errors.New("invalid token")
	}
	user.UserId = Uuid 
	return nil
}

func NewServiceWords(repository repositories.Repository, config *config.Config, logger *zerolog.Logger) *SercviceWords {
	return &SercviceWords{
		repository: repository,
		config: config,
		logger: logger,
	}
}
