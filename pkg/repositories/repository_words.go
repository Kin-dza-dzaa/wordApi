package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

const (
	AddWordQuery 		 = "INSERT INTO user_collection(user_id, word_id, state, collection_name, time_of_last_repeating) VALUES($1, (SELECT id FROM words WHERE word = $2), 0, $3, $4);"
	GetWodsQuery	 	 = "SELECT words.word, user_collection.state, user_collection.collection_name, words.trans_data, user_collection.time_of_last_repeating FROM words INNER JOIN user_collection ON words.id = user_collection.word_id WHERE user_collection.user_id = $1;"
	IfDbHasWordQuery	 = "SELECT EXISTS(SELECT word FROM words WHERE word = $1);"
	IfUserHasWordQuery	 = "SELECT EXISTS(SELECT * FROM user_collection WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));"
	InsertWordDataQuery  = "INSERT INTO words(word, trans_data) VALUES ($1, $2);"
	UpdateWordQuery	 	 = "UPDATE user_collection SET word_id = (SELECT id FROM words WHERE word = $1), state = 0, time_of_last_repeating = $2 WHERE user_id = $3 AND word_id = (SELECT id FROM words WHERE word = $4) AND collection_name = $5;"
	DeleteWordQuery	 	 = "DELETE FROM user_collection WHERE word_id = (SELECT id FROM words WHERE word = $1) AND user_id = $2 AND collection_name = $3;"
	UpdateStateQuery	 = "UPDATE user_collection SET state = $1, time_of_last_repeating = $2 WHERE user_id = $3 AND word_id = (SELECT id FROM words WHERE word = $4) AND collection_name = $5;"
)

type RepositoryWord struct {
	pool *pgxpool.Pool
	config *config.Config
	logger *zerolog.Logger
}

func (repository *RepositoryWord) AddWords(wordsToAdd *models.WordsToAdd, userId string) []string {
	var badWords []string
	words := repository.getAllTranslations(wordsToAdd.Words, &badWords, userId)
	repository.addTranslations(words, &badWords, userId)
    repository.addWordsToUser(&badWords, wordsToAdd, userId)
	return badWords
}

func (repository *RepositoryWord) getAllTranslations(wordsToAdd []models.WordToAdd, badWords *[]string, userId string) []*models.Translation {
	wg := new(sync.WaitGroup)
	channelBadWords := make(chan string, len(wordsToAdd))
	channelJson := make(chan *models.Translation, len(wordsToAdd))
	for _, v := range wordsToAdd {
		exists, err := repository.ifWordInDb(v.Word)
		if err != nil {
			*badWords = append(*badWords, v.Word)
			repository.logger.Warn().Msg(err.Error())
		}
		if !exists {
			wg.Add(1)
			go external.GetTranlations(v.Word, "en", "ru", repository.config, channelJson, channelBadWords, wg, &v)
		}
	}
	wg.Wait()
	close(channelJson)
	close(channelBadWords)
	for v := range channelBadWords {
		*badWords = append(*badWords, v)
	}
	var words []*models.Translation
	for v := range channelJson {
		words = append(words, v)
	}
	return words
}

func (repository *RepositoryWord) ifWordInDb(word string) (bool, error) {
	var result bool
	if err := repository.pool.QueryRow(context.TODO(), IfDbHasWordQuery, word).Scan(&result); err != nil {
		repository.logger.Warn().Msg(err.Error())
		return result, err
	}
	return result, nil
}

func (repository *RepositoryWord) addTranslations(words []*models.Translation, badWords *[]string, userId string) {
	for _, v := range words {
		bytesJson, err := json.Marshal(v)
		if err != nil {
			*badWords = append(*badWords, v.Word)
			repository.logger.Warn().Msg(err.Error())
			continue
		}
		if _, err := repository.pool.Exec(context.TODO(), InsertWordDataQuery, v.Word, string(bytesJson)); err != nil {
			*badWords = append(*badWords, v.Word)
			repository.logger.Warn().Msg(err.Error())
			continue
		}
	}
}

func (repository *RepositoryWord) addWordsToUser(badWords *[]string, wordsToAdd *models.WordsToAdd, userId string) {
	for _, v := range wordsToAdd.Words {
		if !slices.Contains(*badWords, v.Word) {
			if _, err := repository.pool.Exec(context.TODO(), AddWordQuery, userId, v.Word, v.CollectionName, v.TimeOfLastRepeating); err != nil {
				*badWords = append(*badWords, v.Word)
			}
		}
	}
}

func (repository *RepositoryWord) GetWords(userId string) (*models.WordsGet, error) {
	rows, err := repository.pool.Query(context.TODO(), GetWodsQuery, userId)
	if err != nil {
		repository.logger.Error().Msg(err.Error())
		return nil, errors.New("internal error")
	}
	defer rows.Close()
	var WordsGet models.WordsGet
	for rows.Next() {
		var tempWord models.Word
		if err := rows.Scan(&tempWord.Word, &tempWord.State, &tempWord.CollectionName, &tempWord.TransData, &tempWord.TimeOfLastRepeating); err != nil {
			repository.logger.Error().Msg(err.Error())
			return nil, errors.New("internal error")
		}
		WordsGet.Words = append(WordsGet.Words, tempWord)
	}
	return &WordsGet, nil
}

func (repository *RepositoryWord) UpdateWord(wordsToUpdate *models.WordToUpdate, userId string) error {
	var ifUserHasOldWord, ifUserHasNewWord bool
	ifUserHasNewWord, err := repository.ifUserHasWord(wordsToUpdate.NewWord, userId)
	if err != nil {
		return err
	}
	ifUserHasOldWord, err = repository.ifUserHasWord(wordsToUpdate.OldWord, userId)
	if err != nil {
		return err
	}
	if ifUserHasNewWord {
		return fmt.Errorf("you already have word %v", wordsToUpdate.NewWord)
	}
	if ifUserHasOldWord {
		exists, err := repository.ifWordInDb(wordsToUpdate.NewWord)
		if err != nil {
			repository.logger.Error().Msg(err.Error())
			return errors.New("internal error")
		}
		if !exists {
			words := repository.getAllTranslations([]models.WordToAdd{{Word: wordsToUpdate.NewWord, CollectionName: wordsToUpdate.CollectionName, TimeOfLastRepeating: wordsToUpdate.TimeOfLastRepeating}}, &[]string{}, userId)
			repository.addTranslations(words, &[]string{}, userId)
		}
		if response, err := repository.pool.Exec(context.TODO(), UpdateWordQuery, wordsToUpdate.NewWord, wordsToUpdate.TimeOfLastRepeating, userId, wordsToUpdate.OldWord, wordsToUpdate.CollectionName); err != nil || response.RowsAffected() == 0 {
			return errors.New("update was unsuccessful")
		}
	} else {
		return fmt.Errorf("you don't have word %v", wordsToUpdate.OldWord)
	}
	return nil
}

func (repository *RepositoryWord) ifUserHasWord(word string, userId string) (bool, error) {
	var result bool
	err := repository.pool.QueryRow(context.TODO(), IfUserHasWordQuery, userId, word).Scan(&result)
	if err != nil {
		repository.logger.Error().Msg(err.Error())
		return result, errors.New("internal error")
	}
	return result, nil
}

func (repository *RepositoryWord) DeleteWords(words *models.WordsToDelete, userId string) {
	for _, v := range words.Words {
		repository.pool.Exec(context.TODO(), DeleteWordQuery, v.Word, userId, v.CollectionName)
	}
}

func (repository *RepositoryWord) UpdateState(words *models.StatesToUpdate, userId string) {
	for _, v := range words.Words {
		repository.pool.Exec(context.TODO(), UpdateStateQuery, v.NewState, time.Now().UTC(), userId, v.Word, words.CollectionName)
	}
}

func NewRepositoryWord(pool *pgxpool.Pool, logger *zerolog.Logger, config *config.Config) *RepositoryWord {
	return &RepositoryWord{
		pool: pool,
		config: config,
		logger: logger,
	}
}
