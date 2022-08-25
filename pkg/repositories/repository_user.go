package repositories

import (
	"context"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/jackc/pgx/v4"
)

type repositoryUser struct {
	conn *pgx.Conn
}

func (a *repositoryUser) SignUpUser(user *models.User) error {
	sqlCommand := "INSERT INTO USERS(id, user_name, email, password, registration_date) VALUES($1, $2, $3, $4, $5)"
	if _, err := a.conn.Exec(context.TODO(), sqlCommand, user.UserId, user.UserName, user.Email, user.Password, user.Time); err != nil {
		return err
	}
	return nil
}

func (a *repositoryUser) GetUser(user *models.User) (*models.User, error) {
	var tempUser models.User
	if user.Email == "" {
		sqlCommand := "SELECT id, password FROM USERS WHERE user_name = $1"
		row := a.conn.QueryRow(context.TODO(), sqlCommand, user.UserName)
		if err := row.Scan(&tempUser.UserId, &tempUser.Password); err != nil {
			return nil, err
		}
	} else {
		sqlCommand := "SELECT id, password FROM USERS WHERE email = $1"
		row := a.conn.QueryRow(context.TODO(), sqlCommand, user.Email)
		if err := row.Scan(&tempUser.UserId, &tempUser.Password); err != nil {
			return nil, err
		}
	}
	return &tempUser, nil
}
