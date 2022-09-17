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
	Method       	string
	ExpectsError 	bool
	WordsAdd        models.WordsAdd
	WordsUpdate 		models.WordsUpdate
	Result       	[]string
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
		WordsAdd:     models.WordsAdd{Words: [][]string{{"flex", "1", "first"}, {"qweqwe", "1", "first"}, {"ball", "1", "first"}, {"accept", "1", "first"},{ "charge", "1", "first"}, {"battery", "1", "first"}}},
		Result:       []string{"qweqwe"},
	},
	{
		Method:       "AddWords",
		ExpectsError: false,
		WordsAdd:     models.WordsAdd{Words: [][]string{{"", "1", "first"}, {"", "1", "first"}}},
		Result:       []string{"", ""},
	},
	{
		Method:       "AddWords",
		ExpectsError: false,
		WordsAdd:     models.WordsAdd{Words: [][]string{{"as d", "1", "first"}, {"qwewqe", "1", "first"}, {"qwee1e", "1", "first"}, {"qweqw2", "1", "first"}, {"12312wedsa", "1", "first"}, {"asdqwd12", "1", "first"}, {"adsad21", "1", "first"}}},
		Result:       []string{"as d", "qwewqe", "qwee1e", "qweqw2", "12312wedsa", "asdqwd12", "adsad21"},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: false,
		WordsUpdate:  models.WordsUpdate{Words: []string{"flex", "high", "1", "second"}},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: true,
		WordsUpdate:  models.WordsUpdate{Words: []string{"high", "hiasdasdgh", "1", "second"}},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: true,
		WordsUpdate:  models.WordsUpdate{Words: []string{"flex", "trash", "1", "second"}},
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
	if _, err := m.conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1;", testUser.UserId); err != nil{
		m.FailNow(err.Error())
	}
	if _, err := m.conn.Exec(context.Background(), "DELETE FROM words WHERE word=$1;", "high"); err != nil{
		m.FailNow(err.Error())
	}
	if err := m.conn.Close(context.Background()); err != nil {
		m.FailNow(err.Error())
	}
}

func (m *wordSuite) TestAddWord() {
	for _, v := range testsArr {
		if v.Method == "AddWords" {
			res := m.repository.AddWords(v.WordsAdd, testUser.UserId.String())
			m.Equal(res, v.Result)
		}
	}
}

func (m *wordSuite) TestUpdate() {
	for _, v := range testsArr {
		if v.Method == "UpdateWord" {
			res := m.repository.UpdateWord(v.WordsUpdate, testUser.UserId.String())
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
