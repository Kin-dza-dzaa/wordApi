package repositories

import (
	"context"
	"testing"
	"time"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/suite"
)

type UserTest struct {
	models.User
	method   string
	expectedError bool
}

var TestUserData []*UserTest = []*UserTest{
	{
		User: models.User{
			UserId: uuid.New(),
			UserName: "TestName1",
			Email:     "testemail1@gmail.com",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "SignUp",
		expectedError: false,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			UserName: "TestName2",
			Email:     "testemail2@gmail.com",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "SignUp",
		expectedError: false,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			UserName: "TestName1",
			Email:     "testemail1@gmail.com",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "SignUp",
		expectedError: true,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			UserName: "TestName3",
			Email:     "testemail1@gmail.com",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "SignUp",
		expectedError: true,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			UserName: "TestName1",
			Email:     "testemail3@gmail.com",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "SignUp",
		expectedError: true,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			UserName: "TestName1",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "GetUser",
		expectedError: false,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			Email:     "testemail1@gmail.com",
			Password:  "12345",
			Time:      time.Now(),
			},
		method: "GetUser",
		expectedError: false,
	},
	{
		User: models.User{
			UserId: uuid.New(),
			Time:      time.Now(),
			},
		method: "GetUser",
		expectedError: true,
	},
}

type PostgresSuit struct {
	suite.Suite
	repository *Repository
	conn       *pgx.Conn
}

func (p *PostgresSuit) SetupSuite() {
	config, err := config.ReadConfig()
	if err != nil {
		p.FailNow(err.Error())
	}
	db, err := Connect(config.DbUrl)
	if err != nil {
		p.FailNow(err.Error())
	}
	p.conn = db
	p.repository = NewRepository(db, config)
}

func (p *PostgresSuit) TearDownSuite() {
	for _, v := range TestUserData {
		if v.method == "SignUp" && !v.expectedError {
			_, err := p.conn.Exec(context.TODO(), "DELETE FROM USERS WHERE user_name = $1", v.UserName)
			if err != nil {
				p.FailNow(err.Error())
			}
		}
	}
	p.conn.Close(context.TODO())
}

func (p *PostgresSuit) TestRepositorySignUpGetUser() {
	for _, v := range TestUserData {
		switch v.method{
		case "SignUp" : {
			res := p.repository.SignUpUser(&v.User)
			if v.expectedError {
				p.Error(res)
			} else {
				p.Nil(res)
			}
		}
		case "GetUser" : {
			res, err := p.repository.GetUser(&v.User)
			if v.expectedError {
				p.Error(err)
			} else {
				p.Nil(err)
				p.IsType(res, &models.User{})
			}
		}
		}
	}
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(PostgresSuit))
}
