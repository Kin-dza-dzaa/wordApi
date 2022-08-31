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

type userTest struct {
	models.User
	method   string
	expectedError bool
}

var TestUserData []*userTest = []*userTest{
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
			UserName: "TestName2",
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
			Email:     "testemail2@gmail.com",
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

type postgresSuit struct {
	suite.Suite
	repository *Repository
	conn       *pgx.Conn
}

func (p *postgresSuit) SetupSuite() {
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
	err = p.repository.SignUpUser(&models.User{
												UserId: uuid.New(),
												UserName: "TestName1",
												Email:     "testemail1@gmail.com",
												Password:  "12345",
												Time:      time.Now(),
												})
	if err != nil {
		p.FailNow(err.Error())
	}
}

func (p *postgresSuit) TearDownSuite() {
	for _, v := range TestUserData {
		if v.method == "SignUp" && !v.expectedError {
			_, err := p.conn.Exec(context.TODO(), "DELETE FROM USERS WHERE user_name = $1;", v.UserName)
			if err != nil {
				p.FailNow(err.Error())
			}
		}
	}
	_, err := p.conn.Exec(context.TODO(), "DELETE FROM USERS WHERE user_name = $1;", "TestName1")
			if err != nil {
				p.FailNow(err.Error())
			}
	p.conn.Close(context.TODO())
}

func (p *postgresSuit) TestGetUser() {
	for _, v := range TestUserData {
		if v.method == "GetUser" {
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

func (p *postgresSuit) TestSignUp() {
	for _, v := range TestUserData {
		if v.method == "SignUp" {
			res := p.repository.SignUpUser(&v.User)
			if v.expectedError {
				p.Error(res)
			} else {
				p.Nil(res)
			}
		}
	}
}

func TestUserRepo(t *testing.T) {
	suite.Run(t, new(postgresSuit))
}
