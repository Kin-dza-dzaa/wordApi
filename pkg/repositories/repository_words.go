package repositories

import (
	"context"
	"encoding/json"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const (
	queryAddWordToUser  = "INSERT INTO user_collection(user_id, word, state, collection_name, time_of_last_repeating) VALUES($1, $2, 0, $3, $4);"
	queryGetWords       = "SELECT user_collection.word, user_collection.state, user_collection.collection_name, words.trans_data, user_collection.time_of_last_repeating FROM user_collection INNER JOIN words ON words.word = user_collection.word WHERE user_collection.user_id = $1;"
	queryIfDbHasWord    = "SELECT EXISTS(SELECT word FROM words WHERE word = $1);"
	queryIfUserHasWord  = "SELECT EXISTS(SELECT * FROM user_collection WHERE user_id = $1 AND word = $2 AND collection_name = $3);"
	queryInsertWordData = "INSERT INTO words(word, trans_data) VALUES ($1, $2);"
	queryUpdateWord     = "UPDATE user_collection SET word = $1, state = 0, time_of_last_repeating = $2 WHERE user_id = $3 AND word = $4 AND collection_name = $5;"
	queryDeleteWord     = "DELETE FROM user_collection WHERE word = $1 AND user_id = $2 AND collection_name = $3;"
	queryUpdateState    = "UPDATE user_collection SET state = $1, time_of_last_repeating = $2 WHERE user_id = $3 AND word = $4 AND collection_name = $5;"
)

type RepositoryWord struct {
	pool   PgxPool
}

func (repository *RepositoryWord) IfWordInDb(ctx context.Context, word string, result *bool) error {
	if err := repository.pool.QueryRow(ctx, queryIfDbHasWord, word).Scan(result); err != nil {
		return err
	}
	return nil
}

func (repositry *RepositoryWord) IfUserHasWord(ctx context.Context, word string, collectionName string, result *bool, userId string) error {
	if err := repositry.pool.QueryRow(ctx, queryIfUserHasWord, userId, word, collectionName).Scan(result); err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryWord) AddWords(ctx context.Context, words models.WordsToAdd, badwords map[string]string, userId string) error {
	return repository.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, v := range words.Words {
			_, ok := badwords[v.Word]
			if v.TransData != nil && !ok {
				bytesJson, err := json.Marshal(v.TransData)
				if err != nil {
					return err
				}
				if _, err := tx.Exec(ctx, queryInsertWordData, v.Word, string(bytesJson)); err != nil {
					return err
				}
				if _, err := tx.Exec(ctx, queryAddWordToUser, userId, v.Word, v.CollectionName, v.TimeOfLastRepeating); err != nil {
					return err
				}
			} else {
				if !ok {
					if _, err := tx.Exec(ctx, queryAddWordToUser, userId, v.Word, v.CollectionName, v.TimeOfLastRepeating); err != nil {
						continue
					}
				}
			}
		}
		return nil
	})
}

func (repository *RepositoryWord) GetWords(ctx context.Context, words models.Words, userId string) error {
	rows, err := repository.pool.Query(ctx, queryGetWords, userId)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var tempWord models.Word
		if err := rows.Scan(&tempWord.Word, &tempWord.State, &tempWord.CollectionName, &tempWord.TransData, &tempWord.TimeOfLastRepeating); err != nil {
			return err
		}
		*words.Words = append((*words.Words), tempWord)
	}
	return nil
}

func (repository *RepositoryWord) UpdateWord(ctx context.Context, wordToUpdate models.WordToUpdate, userId string) error {
	return repository.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if wordToUpdate.TransData != nil {
			bytesJson, err := json.Marshal(wordToUpdate.TransData)
			if err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, queryInsertWordData, wordToUpdate.NewWord, string(bytesJson)); err != nil {
				return err
			}
		}

		_, err := tx.Exec(ctx, queryUpdateWord, wordToUpdate.NewWord, wordToUpdate.TimeOfLastRepeating, userId, wordToUpdate.OldWord, wordToUpdate.CollectionName)
		if err != nil {
			return err
		}

		return nil
	})
}

func (repository *RepositoryWord) DeleteWords(ctx context.Context, words models.WordsToDelete, userId string) error {
	return repository.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, v := range words.Words {
			if _, err := tx.Exec(ctx, queryDeleteWord, v.Word, userId, v.CollectionName); err != nil {
				return err
			}
		}
		return nil
	})
}

func (repository *RepositoryWord) UpdateState(ctx context.Context, words models.StatesToUpdate, userId string) error {
	return repository.pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, v := range words.Words {
			if _, err := tx.Exec(ctx, queryUpdateState, v.NewState, v.TimeOfLastRepeating, userId, v.Word, words.CollectionName); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewRepositoryWord(pool PgxPool) *RepositoryWord {
	return &RepositoryWord{
		pool:   pool,
	}
}

type PgxPool interface {
	Close()
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
}
