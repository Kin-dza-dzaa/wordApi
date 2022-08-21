package repository

import (
	"context"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4"
)

type repository struct {
	conn *pgx.Conn
}

func (a *repository) SignUpUser(user *models.User) error {
	sqlCommand := "INSERT INTO USERS(id, user_name, email, password, registration_date) VALUES($1, $2, $3, $4, $5)"
	if _, err := a.conn.Exec(context.Background(), sqlCommand, user.User_id, user.User_name, user.Email, user.Password, user.Time); err != nil {
		return err
	}
	return nil
}

func (a *repository) GetUser(email string) (*models.User, error) {
	var tempUser models.User
	sqlCommand := "SELECT id, password FROM USERS WHERE email = $1"
	row := a.conn.QueryRow(context.Background(), sqlCommand, email)
	if err := row.Scan(tempUser.User_id, tempUser.Password); err != nil {
		return nil, err
	}
	return &tempUser, nil
}
