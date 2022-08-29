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

func (r *repositoryWord) AddWords(words models.Words, userId string) []string {
	var badWords []string
	for _, v := range words.Words{
		var ifDBHasWord bool
		err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT word FROM words WHERE word = $1);", v).Scan(&ifDBHasWord)
		if err != nil {
			badWords = append(badWords, v)
			continue
		}
		if !ifDBHasWord {
			translationJson, err :=  external.GetTranlations(v, "en", "ru", r.config)
			if err != nil {
				badWords = append(badWords, v)
				continue
			}
			stringJson, err := json.Marshal(translationJson)
			if err != nil {
				badWords = append(badWords, v)
				continue
			}
			commandTag, _ := r.conn.Exec(context.TODO(), "INSERT INTO words(word, trans_data) VALUES ($1, $2);", v, string(stringJson))
			if commandTag.RowsAffected() == 0 {
				badWords = append(badWords, v)
				continue
			}
		}
		r.conn.Exec(context.TODO(), "INSERT INTO userword(user_id, word_id) VALUES($1, (SELECT id FROM words WHERE word = $2));", userId, v)
	}
	return badWords
}

func (r *repositoryWord) GetWords(userId string) (*[]external.Translation, error) {
	sql := `
			SELECT trans_data FROM words
				WHERE id IN (SELECT word_id FROM userword WHERE user_id = $1);
		   ` 
	rows, err := r.conn.Query(context.TODO(), sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var wordStruct []external.Translation
	for rows.Next() {
		var tempWord external.Translation
		if err := rows.Scan(&tempWord); err != nil {
			return nil, err
		}
		wordStruct = append(wordStruct, tempWord)
	}
	return &wordStruct, nil
}

func (r *repositoryWord) UpdateWord(words models.Words, userId string) error {
	if len(words.Words) != 2 || words.Words[0] == words.Words[1] {
		return fmt.Errorf("expected 2 not even words got %d", len(words.Words))
	}
	var ifUserHasOldWord, ifUserHasNewWord bool
	err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM userword WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.Words[0]).Scan(&ifUserHasOldWord)
	if err != nil {
		return err
	}
	if !ifUserHasOldWord {
		return fmt.Errorf("you don't have %s word", words.Words[0])
	}
	err = r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM userword WHERE user_id = $1 AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.Words[1]).Scan(&ifUserHasNewWord)
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
	if response, err := r.conn.Exec(context.TODO(), "UPDATE userword SET word_id = (SELECT id FROM words WHERE word = $1) WHERE user_id = $2 AND word_id = (SELECT id FROM words WHERE word = $3);", words.Words[1], userId, words.Words[0]); err != nil || response.RowsAffected() == 0 {
		return errors.New("update was unsuccessful")
	}
	return nil
}

func (r *repositoryWord) DeleteWords(words models.Words, userId string) {
	for _, v := range words.Words {
		r.conn.Exec(context.TODO(), "DELETE FROM userword WHERE word_id = (SELECT id FROM words WHERE word = $1) AND user_id = $2;", v, userId)
	}
}
