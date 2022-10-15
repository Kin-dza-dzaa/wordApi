package repositories

// this test uses real db connection, keep it in mind

import (
	"context"
	"os"
	"testing"
	"time"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type tests struct {
	Method       string
	ExpectsError bool
	WordsToAdd   models.WordsToAdd
	WordsUpdate  models.WordToUpdate
	Result       []string
	err string
}

type wordSuite struct {
	suite.Suite
	UUidTest 	uuid.UUID
	pool        *pgxpool.Pool
	repository	Repository
}

func (suite *wordSuite) SetupSuite() {
	suite.UUidTest = uuid.New()
	config, err := config.ReadConfig()
	if err != nil {
		suite.FailNow(err.Error())
	}
	conn, err := pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.pool = conn
	myLogger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	suite.repository = NewRepository(conn, &myLogger, config)
}

func (suite *wordSuite) TearDownSuite() {
	var wordsToDelte []string = []string{"flex", "high", "ball", "accept", "charge", "battery"}
	for _, v := range wordsToDelte {
		if _, err := suite.pool.Exec(context.TODO(), "DELETE FROM user_collection WHERE word_id=(select id from words where word = $1);", v); err != nil{
			suite.FailNow(err.Error())
		}
		if _, err := suite.pool.Exec(context.TODO(), "DELETE FROM words WHERE word = $1;", v); err != nil{
			suite.FailNow(err.Error())
		}
	}
	suite.pool.Close()
}

func (suite *wordSuite) TestAddWord() {
	var testSlice []tests = []tests{
		{
			Method:       "AddWords",
			ExpectsError: false,
			WordsToAdd:     models.WordsToAdd{Words: []models.WordToAdd{{Word: "flex", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "qweqwe", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "ball", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "accept", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "charge", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "battery", CollectionName: "first", TimeOfLastRepeating: time.Now()}}},
			Result:       []string{"qweqwe"},
		},
		{
			Method:       "AddWords",
			ExpectsError: false,
			WordsToAdd:     models.WordsToAdd{Words: []models.WordToAdd{{Word: "", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "", CollectionName: "first", TimeOfLastRepeating: time.Now()}}},
			Result:       []string{"", ""},
		},
		{
			Method:       "AddWords",
			ExpectsError: false,
			WordsToAdd:     models.WordsToAdd{Words: []models.WordToAdd{{Word: "as d", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "qwewqe", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "qwee1e", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "qweqw2", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "12312wedsa", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "asdqwd12", CollectionName: "first", TimeOfLastRepeating: time.Now()}, {Word: "adsad21", CollectionName: "first", TimeOfLastRepeating: time.Now()}}},
			Result:       []string{"as d", "qwewqe", "qwee1e", "qweqw2", "12312wedsa", "asdqwd12", "adsad21"},
		},
	} 
	for _, v := range testSlice {
		res := suite.repository.AddWords(&v.WordsToAdd, suite.UUidTest.String())
		for _, value := range res {
			suite.Contains(v.Result, value)
		}
	}
}

func (suite *wordSuite) TestUpdate() {
	var testSlice []tests = []tests{
		{
			Method:       "UpdateWord",
			ExpectsError: false,
			WordsUpdate:  models.WordToUpdate{OldWord: "flex", NewWord: "high", CollectionName: "first", TimeOfLastRepeating: time.Now()},
		},
		{
			Method:       "UpdateWord",
			ExpectsError: true,
			WordsUpdate:  models.WordToUpdate{OldWord: "high", NewWord: "hiasdasdgh", CollectionName: "second", TimeOfLastRepeating: time.Now()},
			err: "update was unsuccessful",
		},
		{
			Method:       "UpdateWord",
			ExpectsError: true,
			WordsUpdate:  models.WordToUpdate{OldWord: "flex", NewWord: "trash", CollectionName: "second", TimeOfLastRepeating: time.Now()},
			err: "you don't have word flex",
		},
	} 
	for _, v := range testSlice {
		res := suite.repository.UpdateWord(&v.WordsUpdate, suite.UUidTest.String())
		if v.ExpectsError {
			suite.EqualError(res, v.err)
		} else {
			suite.Nil(res)
		}
	}
}

func (suite *wordSuite) TestGetWords() {
	var testSlice []tests = []tests{
		{
			Method:       "GetWords",
			ExpectsError: false,
		},
		{
			Method:       "GetWords",
			ExpectsError: false,
		},
	} 
	for _, v := range testSlice {
		if v.ExpectsError {
			_, err := suite.repository.GetWords(uuid.New().String())
			suite.EqualError(err, v.err)
		} else {
			_, err := suite.repository.GetWords(suite.UUidTest.String())
			suite.Nil(err)
		}
	}
}

func TestWordRepo(t *testing.T) {
	suite.Run(t, new(wordSuite))
}
