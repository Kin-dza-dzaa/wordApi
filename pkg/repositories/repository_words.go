package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4"
)

type repositoryWord struct {
	conn *pgx.Conn
	config *config.Config
}

func (r *repositoryWord) AddWords(words models.WordsAdd, userId string) []string {
	var badWords []string
	sql := `INSERT INTO user_collection(user_id, word_id, state, collection_name) 
				VALUES($1, (SELECT id FROM words WHERE word = $2), 1, $3);`
	for _, v := range words.Words{
		var ifDBHasWord bool
		err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT word FROM words WHERE word = $1);", v.Word).Scan(&ifDBHasWord)
		if err != nil {
			badWords = append(badWords, v.Word)
			continue
		}
		if !ifDBHasWord {
			translationJson, err :=  external.GetTranlations(v.Word, "en", "ru", r.config)
			if err != nil {
				badWords = append(badWords, v.Word)
				continue
			}
			bytesJson, err := json.Marshal(translationJson)
			if err != nil {
				badWords = append(badWords, v.Word)
				continue
			}
			r.conn.Exec(context.TODO(), "INSERT INTO words(word, trans_data) VALUES ($1, $2);", v.Word, string(bytesJson))
		}
		if commandTag, err := r.conn.Exec(context.TODO(), sql, userId, v.Word, v.CollectionName); commandTag.RowsAffected() == 0 || err != nil {
			badWords = append(badWords, v.Word)
		}	
	}
	return badWords
}

func (r *repositoryWord) GetWords(userId string) (*models.WordsGet, error) {
	sql := `
			SELECT words.word, user_collection.state, user_collection.collection_name, words.trans_data 
				FROM words
				INNER JOIN user_collection 
					ON (words.id = user_collection.word_id) 
						WHERE (user_collection.user_id = $1);
		   ` 
	rows, err := r.conn.Query(context.TODO(), sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var wordsStruct models.WordsGet
	for rows.Next() {
		var tempWord models.Word
		if err := rows.Scan(&tempWord.Word, &tempWord.State, &tempWord.CollectionName, &tempWord.TransData); err != nil {
			return nil, err
		}
		wordsStruct.Words = append(wordsStruct.Words, tempWord)
	}
	return &wordsStruct, nil
}

func (r *repositoryWord) UpdateWord(words models.WordToUpdate, userId string) error {
	var ifUserHasOldWord, ifUserHasNewWord bool
	err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM user_collection WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.OldWord).Scan(&ifUserHasOldWord)
	if err != nil {
		return err
	}
	if !ifUserHasOldWord {
		return fmt.Errorf("you don't have %s word", words.OldWord)
	}
	err = r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM user_collection WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.NewWord).Scan(&ifUserHasNewWord)
	if err != nil {
		return err
	}
	if !ifUserHasNewWord {
		translationJson, err :=  external.GetTranlations(words.NewWord, "en", "ru", r.config)
		if err != nil {
			return err
		}
		jsonString, err := json.Marshal(translationJson)
		if err != nil {
			return fmt.Errorf("something wrong with word %s", words.NewWord)
		}
		r.conn.Exec(context.TODO(), "INSERT INTO words(word, trans_data) VALUES($1, $2);", words.NewWord, string(jsonString))
	}
	if response, err := r.conn.Exec(context.TODO(), "UPDATE user_collection SET word_id = (SELECT id FROM words WHERE word = $1), state = 1, collection_name = $2 WHERE user_id = $3 AND word_id = (SELECT id FROM words WHERE word = $4);", words.NewWord, words.CollectionName, userId, words.OldWord); err != nil || response.RowsAffected() == 0 {
		return errors.New("update was unsuccessful")
	}
	return nil
}

func (r *repositoryWord) DeleteWords(words models.WordsDelete, userId string) {
	for _, v := range words.Words {
		r.conn.Exec(context.TODO(), "DELETE FROM user_collection WHERE word_id = (SELECT id FROM words WHERE word = $1) AND user_id = $2 AND collection_name = $3;", v[0], userId, v[1])
	}
}
