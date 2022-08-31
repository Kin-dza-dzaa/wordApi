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

type tests struct {
	Method       string
	ExpectsError bool
	Words        models.Words
	Result       []string
}

var testUser models.User = models.User{
	UserId:   uuid.New(),
	UserName: "TestUserWords",
	Email:    "TestUserWords@gmail.com",
	Password: "1312e9djc091j20e91092w",
	Time:     time.Now(),
}

var testsArr []tests = []tests{
	{
		Method:       "AddWords",
		ExpectsError: false,
		Words:        models.Words{Words: []string{"flex", "flex", "qweqwe", "ball", "accept", "charge", "battery"}},
		Result:       []string{"qweqwe"},
	},
	{
		Method:       "AddWords",
		ExpectsError: false,
		Words:        models.Words{Words: []string{"", ""}},
		Result:       []string{"", ""},
	},
	{
		Method:       "AddWords",
		ExpectsError: false,
		Words:        models.Words{Words: []string{"asd", "qwewqe", "qwee1e", "qweqw2", "12312wedsa", "asdqwd12", "adsad21"}},
		Result:       []string{"asd", "qwewqe", "qwee1e", "qweqw2", "12312wedsa", "asdqwd12", "adsad21"},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: false,
		Words:        models.Words{Words: []string{"flex", "high"}},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: true,
		Words:        models.Words{Words: []string{"high", "asdasdasdasd"}},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: true,
		Words:        models.Words{Words: []string{"flex", "true"}},
	},
}

type wordSuite struct {
	suite.Suite
	conn       *pgx.Conn
	repository *Repository
}

func (m *wordSuite) SetupSuite() {
	config, err := config.ReadConfig()
	if err != nil {
		m.FailNow(err.Error())
	}
	conn, err := Connect(config.DbUrl)
	if err != nil {
		m.FailNow(err.Error())
	}
	m.conn = conn
	m.repository = NewRepository(conn, config)
	err = m.repository.SignUpUser(&testUser)
	if err != nil {
		m.FailNow(err.Error())
	}
}

func (m *wordSuite) TearDownSuite() {
	_, err := m.conn.Exec(context.Background(), "DELETE FROM USERS WHERE id=$1;", testUser.UserId)
	if err != nil {
		m.FailNow(err.Error())
	}
	err = m.conn.Close(context.Background())
	if err != nil {
		m.FailNow(err.Error())
	}
}

func (m *wordSuite) TestAddWord() {
	for _, v := range testsArr {
		if v.Method == "AddWords" {
			res := m.repository.AddWords(v.Words, testUser.UserId.String())
			m.Equal(res, v.Result)
		}
	}
}

func (m *wordSuite) TestUpdate() {
	for _, v := range testsArr {
		if v.Method == "UpdateWord" {
			res := m.repository.UpdateWord(v.Words, testUser.UserId.String())
			if v.ExpectsError {
				m.Error(res)
			} else {
				m.Nil(res)
			}
		}
	}
}

func TestWordRepo(t *testing.T) {
	suite.Run(t, new(wordSuite))
}
