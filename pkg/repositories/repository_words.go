package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4"
)

type repositoryWord struct {
	conn *pgx.Conn
}

func (r *repositoryWord) AddWords(words models.Words, userId string) {
	for _, v := range words.Words{
		r.conn.Exec(context.TODO(), "INSERT INTO words(word) VALUES ($1);", v)
		r.conn.Exec(context.TODO(), "INSERT INTO userword(user_id, word_id) VALUES(CAST($1 AS UUID), (SELECT id FROM words WHERE word = $2));", userId, v)
	}
}

func (r *repositoryWord) GetWords(userId string) (*models.Words, error) {
	sql := `
			SELECT word FROM words 
				WHERE id IN (SELECT word_id FROM userword WHERE user_id = CAST($1 AS UUID));
		   ` 
	rows, err := r.conn.Query(context.TODO(), sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var wordStruct *models.Words = new(models.Words)
	for rows.Next() {
		var tempWord string
		if err := rows.Scan(&tempWord); err != nil {
			return nil, err
		}
		wordStruct.Words = append(wordStruct.Words, tempWord)
	}
	return wordStruct, nil
}

func (r *repositoryWord) UpdateWord(words models.Words, userId string) error {
	if len(words.Words) != 2 || words.Words[0] == words.Words[1] {
		return fmt.Errorf("expected 2 not even words got %d", len(words.Words))
	}
	var ifUserHasOldWord bool
	err := r.conn.QueryRow(context.TODO(), "SELECT EXISTS(SELECT * FROM userword WHERE user_id = CAST($1 AS UUID) AND word_id = (SELECT id FROM words WHERE word = $2));", userId, words.Words[0]).Scan(&ifUserHasOldWord)
	if err != nil {
		return err
	}
	if !ifUserHasOldWord {
		return fmt.Errorf("you don't have %s word", words.Words[0])
	}
	r.conn.Exec(context.TODO(), "INSERT INTO words(word) VALUES($1);", words.Words[1])
	if response, err := r.conn.Exec(context.TODO(), "UPDATE userword SET word_id = (SELECT id FROM words WHERE word = $1) WHERE user_id = CAST($2 AS UUID) AND word_id = (SELECT id FROM words WHERE word = $3);", words.Words[1], userId, words.Words[0]); err != nil || response.RowsAffected() == 0 {
		return errors.New("update was unsuccessful")
	}
	return nil
}

func (r *repositoryWord) DeleteWords(words models.Words, userId string) {
	for _, v := range words.Words {
		r.conn.Exec(context.TODO(), "DELETE FROM userword WHERE word_id = (SELECT id FROM words WHERE word = $1);", v)
	}
}
