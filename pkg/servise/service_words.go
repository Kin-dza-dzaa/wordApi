package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/apierror"
	external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyHasNewWord = errors.New("you already has new word")
	ErrUserNotHaveOldWord    = errors.New("you don't have old word")
	ErrInvalidToken          = errors.New("invalid token")
)

type SercviceWords struct {
	repository repositories.Repository
	config     *config.Config
}

func (service *SercviceWords) AddWords(ctx context.Context, words models.WordsToAdd, userId string) (map[string]string, error) {
	badWords := make(map[string]string, len(words.Words))
	wg := new(sync.WaitGroup)
	channelBadWords := make(chan string, len(words.Words))

	for i := range words.Words {
		words.Words[i].TimeOfLastRepeating = time.Now().UTC()
		var exists, hasWord bool
		if err := service.repository.IfWordInDb(ctx, words.Words[i].Word, &exists); err != nil {
			return badWords, err
		}
		if err := service.repository.IfUserHasWord(ctx, words.Words[i].Word, words.Words[i].CollectionName, &hasWord, userId); err != nil {
			return badWords, err
		}
		if !exists {
			words.Words[i].TransData = new(models.Translation)
			wg.Add(1)
			go external.GetTranlations(words.Words[i].Word, words.Words[i].TransData, "en", "ru", service.config, channelBadWords, wg)
		}
		if hasWord {
			badWords[words.Words[i].Word] = words.Words[i].Word
		}
	}

	wg.Wait()
	close(channelBadWords)

	for v := range channelBadWords {
		badWords[v] = v
	}

	if err := service.repository.AddWords(ctx, words, badWords, userId); err != nil {
		return badWords, err
	}

	return badWords, nil
}

func (service *SercviceWords) UpdateWord(ctx context.Context, words models.WordToUpdate, userId string) error {
	words.TimeOfLastRepeating = time.Now().UTC()
	var ifUserHasOldWord, ifUserHasNewWord, ifNewWordInDb bool

	if err := service.repository.IfUserHasWord(ctx, words.NewWord, words.CollectionName, &ifUserHasNewWord, userId); err != nil {
		return err
	}

	if err := service.repository.IfUserHasWord(ctx, words.OldWord, words.CollectionName, &ifUserHasOldWord, userId); err != nil {
		return err
	}

	if err := service.repository.IfWordInDb(ctx, words.NewWord, &ifNewWordInDb); err != nil {
		return err
	}

	if ifUserHasNewWord {
		return apierror.NewResponse("error", ErrUserAlreadyHasNewWord.Error(), http.StatusBadRequest)
	}

	if !ifUserHasOldWord {
		return apierror.NewResponse("error", ErrUserNotHaveOldWord.Error(), http.StatusBadRequest)
	}

	if !ifNewWordInDb {
		words.TransData = new(models.Translation)
		if err := external.GetTranlationsUpdate(words.NewWord, words.TransData, "en", "ru", service.config); err != nil {
			return err
		}
	}
	return service.repository.UpdateWord(ctx, words, userId)
}

func (service *SercviceWords) GetWords(ctx context.Context, words models.Words, userId string) error {
	return service.repository.GetWords(ctx, words, userId)
}

func (service *SercviceWords) DeleteWords(ctx context.Context, words models.WordsToDelete, userId string) error {
	return service.repository.DeleteWords(ctx, words, userId)
}

func (service *SercviceWords) UpdateState(ctx context.Context, words models.StatesToUpdate, userId string) error {
	for i := range words.Words {
		words.Words[i].TimeOfLastRepeating = time.Now().UTC()
	}
	return service.repository.UpdateState(ctx, words, userId)
}

func (service *SercviceWords) ValidateToken(user *models.User) error {
	token, err := jwt.ParseWithClaims(user.Jwt, &models.MyJwtClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(service.config.JWTString), nil })
	if err != nil {
		return apierror.NewResponse("error", ErrInvalidToken.Error(), http.StatusUnauthorized)
	}
	claims, ok := token.Claims.(*models.MyJwtClaims)
	if !ok || !token.Valid {
		return apierror.NewResponse("error", ErrInvalidToken.Error(), http.StatusUnauthorized)
	}
	user.CsrfToken = claims.XCSRFToken
	Uuid, err := uuid.Parse(claims.UserId)
	if err != nil {
		return apierror.NewResponse("error", ErrInvalidToken.Error(), http.StatusUnauthorized)
	}
	user.UserId = Uuid
	return nil
}

func NewServiceWords(repository repositories.Repository, config *config.Config) *SercviceWords {
	return &SercviceWords{
		repository: repository,
		config:     config,
	}
}
