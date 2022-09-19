package repositories

// this test uses real db connection, keep it in mind

import (
	"context"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type tests struct {
	Method       string
	ExpectsError bool
	WordsAdd     models.WordsAdd
	WordsUpdate  models.WordToUpdate
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
		WordsAdd:     models.WordsAdd{Words: []models.WordToAdd{{Word: "flex", CollectionName: "first"}, {Word: "qweqwe", CollectionName: "first"}, {Word: "ball", CollectionName: "first"}, {Word: "accept", CollectionName: "first"}, {Word: "charge", CollectionName: "first"}, {Word: "battery", CollectionName: "first"}}},
		Result:       []string{"qweqwe"},
	},
	{
		Method:       "AddWords",
		ExpectsError: false,
		WordsAdd:     models.WordsAdd{Words: []models.WordToAdd{{Word: "", CollectionName: "first"}, {Word: "", CollectionName: "first"}}},
		Result:       []string{"", ""},
	},
	{
		Method:       "AddWords",
		ExpectsError: false,
		WordsAdd:     models.WordsAdd{Words: []models.WordToAdd{{Word: "as d", CollectionName: "first"}, {Word: "qwewqe", CollectionName: "first"}, {Word: "qwee1e", CollectionName: "first"}, {Word: "qweqw2", CollectionName: "first"}, {Word: "12312wedsa", CollectionName: "first"}, {Word: "asdqwd12", CollectionName: "first"}, {Word: "adsad21", CollectionName: "first"}}},
		Result:       []string{"as d", "qwewqe", "qwee1e", "qweqw2", "12312wedsa", "asdqwd12", "adsad21"},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: false,
		WordsUpdate:  models.WordToUpdate{OldWord: "flex", NewWord: "high", CollectionName: "second"},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: true,
		WordsUpdate:  models.WordToUpdate{OldWord: "high", NewWord: "hiasdasdgh", CollectionName: "second"},
	},
	{
		Method:       "UpdateWord",
		ExpectsError: true,
		WordsUpdate:  models.WordToUpdate{OldWord: "flex", NewWord: "trash", CollectionName: "second"},
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
		var wordsToDelte []string = []string{"flex", "high", "ball", "accept", "charge", "battery"}
	if _, err := m.conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1;", testUser.UserId); err != nil {
		m.FailNow(err.Error())
	}
	for _, v := range wordsToDelte {
		if _, err := m.conn.Exec(context.Background(), "DELETE FROM words WHERE word=$1;", v); err != nil{
			m.FailNow(err.Error())
		}
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
