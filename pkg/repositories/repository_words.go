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
				VALUES($1, (SELECT id FROM words WHERE word = $2), $3, $4);`
	for _, v := range words.Words{
		var ifDBHasWord bool
		err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT word FROM words WHERE word = $1);", v[0]).Scan(&ifDBHasWord)
		if err != nil {
			badWords = append(badWords, v[0])
			continue
		}
		if !ifDBHasWord {
			translationJson, err :=  external.GetTranlations(v[0], "en", "ru", r.config)
			if err != nil {
				badWords = append(badWords, v[0])
				continue
			}
			bytesJson, err := json.Marshal(translationJson)
			if err != nil {
				badWords = append(badWords, v[0])
				continue
			}
			r.conn.Exec(context.TODO(), "INSERT INTO words(word, trans_data) VALUES ($1, $2);", v[0], string(bytesJson))
		}
		if commandTag, err := r.conn.Exec(context.TODO(), sql, userId, v[0], v[1], v[2]); commandTag.RowsAffected() == 0 || err != nil {
			badWords = append(badWords, v[0])
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

func (r *repositoryWord) UpdateWord(words models.WordsUpdate, userId string) error {
	var ifUserHasOldWord, ifUserHasNewWord bool
	err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM user_collection WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.Words[0]).Scan(&ifUserHasOldWord)
	if err != nil {
		return err
	}
	if !ifUserHasOldWord {
		return fmt.Errorf("you don't have %s word", words.Words[0])
	}
	err = r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM user_collection WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.Words[1]).Scan(&ifUserHasNewWord)
	if err != nil {
		return err
	}
	if !ifUserHasNewWord {
		translationJson, err :=  external.GetTranlations(words.Words[1], "en", "ru", r.config)
		if err != nil {
			return err
		}
		jsonString, err := json.Marshal(translationJson)
		if err != nil {
			return fmt.Errorf("something wrong with word %s", words.Words[1])
		}
		r.conn.Exec(context.TODO(), "INSERT INTO words(word, trans_data) VALUES($1, $2);", words.Words[1], string(jsonString))
	}
	if response, err := r.conn.Exec(context.TODO(), "UPDATE user_collection SET word_id = (SELECT id FROM words WHERE word = $1), state = $2, collection_name = $3 WHERE user_id = $4 AND word_id = (SELECT id FROM words WHERE word = $5);", words.Words[1], words.Words[2], words.Words[3], userId, words.Words[0]); err != nil || response.RowsAffected() == 0 {
		return errors.New("update was unsuccessful")
	}
	return nil
}

func (r *repositoryWord) DeleteWords(words models.WordsDelete, userId string) {
	for _, v := range words.Words {
		r.conn.Exec(context.TODO(), "DELETE FROM user_collection WHERE word_id = (SELECT id FROM words WHERE word = $1) AND user_id = $2 AND collection_name = $3;", v[0], userId, v[1])
	}
}
