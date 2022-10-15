package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/handlers"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/wordApi/pkg/servise"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type testStruct struct {
	message string 
	wrongCsrf bool
	body string
}

type Response struct {
	Message string		`json:"message,omitempty"`
	Result string		`json:"result"`
}

type IntegrationTestSuite struct {
	suite.Suite
	server *http.Server
	pool *pgxpool.Pool
	jwt string
	CsrfToken string
}

func (suite *IntegrationTestSuite) SetupSuite() {
	myLogger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	config, err := config.ReadConfig()
	if err != nil {
		suite.FailNow(err.Error())
	}
	pool, err := pgxpool.Connect(context.TODO(), config.DbUrl)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.pool = pool
	myRepository := repositories.NewRepository(pool, &myLogger, config)
	myService := service.NewService(myRepository, config, &myLogger)
	myHandlers := handlers.NewHandlers(myService)
	myHandlers.InitilizeHandlers()
	srv := &http.Server{
		Handler: myHandlers.Cors.Handler(myHandlers.Router),
		Addr:    config.Adress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	suite.server = srv
	go func() {
		myLogger.Info().Msg(fmt.Sprintf("Staring server wordapi at %v", config.Adress))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			suite.FailNow(err.Error())
		}
	}()
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, models.MyJwtClaims{
		UserId: uuid.New().String(),
		XCSRFToken: "csrf token",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(5*time.Minute).Unix(),
		},
	})
	token, err := jwt.SignedString([]byte(config.JWTString))
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.jwt = token
	suite.CsrfToken = "csrf token"
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	var wordsToDelte []string = []string{"test", "go", "new"}
	for _, v := range wordsToDelte {
		if _, err := suite.pool.Exec(context.TODO(), "DELETE FROM user_collection WHERE word_id=(select id from words where word = $1);", v); err != nil{
			suite.FailNow(err.Error())
		}
		if _, err := suite.pool.Exec(context.TODO(), "DELETE FROM words WHERE word = $1;", v); err != nil{
			suite.FailNow(err.Error())
		}
	}
	suite.pool.Close()
	if err := suite.server.Shutdown(context.TODO()); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *IntegrationTestSuite) TestGetWords() {
	var testSlice []testStruct = []testStruct{
		{
			wrongCsrf: false,
			message: "words were sent",
		},
		{
			wrongCsrf: true,
			message: "invalid token",
		},
	}
	for _, v := range testSlice {
		var result Response
		if v.wrongCsrf {
			r, err := http.NewRequest("GET", "http://localhost:8000/words", nil)
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndBadCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		} else {
			r, err := http.NewRequest("GET", "http://localhost:8000/words", nil)
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *IntegrationTestSuite) TestAddWords() {
	var testSlice []testStruct = []testStruct{
		{
			wrongCsrf: false,
			message: "some words weren't added",
			body: `{"words": [{"word": "test", "collection_name": "test"}, {"word": "go", "collection_name": "test"}, {"word": "asdasd", "collection_name": "test"}]}`,
		},
		{
			wrongCsrf: true,
			message: "invalid token",
			body: `{"words": []}`,
		},
	}
	for _, v := range testSlice {
		var result Response
		if v.wrongCsrf {
			r, err := http.NewRequest("POST", "http://localhost:8000/words", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndBadCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		} else {
			r, err := http.NewRequest("POST", "http://localhost:8000/words", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *IntegrationTestSuite) TestUpdateWord() {
	var testSlice []testStruct = []testStruct{
		{
			wrongCsrf: false,
			message: `you don't have word `,
			body: `{}`,
		},
		{
			wrongCsrf: true,
			message: `invalid token`,
			body: `{}`,
		},
		{
			wrongCsrf: false,
			message: `you already have word test`,
			body: `{"old_word": "go", "new_word": "test", "collection_name": "test"}`,
		},
		{
			wrongCsrf: false,
			message: `update was unsuccessful`,
			body: `{"old_word": "go", "new_word": "wrong_input", "collection_name": "test"}`,
		},
		{
			wrongCsrf: false,
			message: `update was unsuccessful`,
			body: `{"old_word": "go", "new_word": "new", "collection_name": "test1"}`,
		},
	}
	for _, v := range testSlice {
		var result Response
		if v.wrongCsrf {
			r, err := http.NewRequest("PUT", "http://localhost:8000/words", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndBadCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		} else {
			r, err := http.NewRequest("PUT", "http://localhost:8000/words", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *IntegrationTestSuite) TestUpdateStates() {
	var testSlice []testStruct = []testStruct{
		{
			wrongCsrf: false,
			message: `expected object: {"collection_name": string, "words": [{"word": string, "new_state": number}]}`,
			body: `{}`,
		},
		{
			wrongCsrf: false,
			message: `words were updated`,
			body: `{"collection_name": "test1", "words": []}`,
		},
		{
			wrongCsrf: true,
			message: `invalid token`,
			body: `{}`,
		},
	}
	for _, v := range testSlice {
		var result Response
		if v.wrongCsrf {
			r, err := http.NewRequest("PUT", "http://localhost:8000/words/state", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndBadCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		} else {
			r, err := http.NewRequest("PUT", "http://localhost:8000/words/state", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *IntegrationTestSuite) TestZDeleteWords() {
	var testSlice []testStruct = []testStruct{
		{
			wrongCsrf: false,
			message: `expected object: {"words": [{"collection_name": string, "word": string}...]}`,
			body: `{}`,
		},
		{
			wrongCsrf: false,
			message: `words were delted`,
			body: `{"words": [{"collection_name": "test", "word": "go"}, {"collection_name": "test", "word": "test"}, {"collection_name": "test", "word": "asdasdasd"}]}`,
		},
		{
			wrongCsrf: true,
			message: `invalid token`,
			body: `{}`,
		},
	}
	for _, v := range testSlice {
		var result Response
		if v.wrongCsrf {
			r, err := http.NewRequest("DELETE", "http://localhost:8000/words", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndBadCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		} else {
			r, err := http.NewRequest("DELETE", "http://localhost:8000/words", bytes.NewReader([]byte(v.body)))
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.setCookieAndCsrf(r)
			suite.sendRequest(r, &result)
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *IntegrationTestSuite) setCookieAndCsrf(r *http.Request) {
	r.Header.Set("X-CSRF-Token", suite.CsrfToken)
	r.AddCookie(&http.Cookie{
		Name: "Access-token",
		Value: suite.jwt,
	})
}

func (suite *IntegrationTestSuite) setCookieAndBadCsrf(r *http.Request) {
	r.Header.Set("X-CSRF-Token", "Bad csrf")
	r.AddCookie(&http.Cookie{
		Name: "Access-token",
		Value: suite.jwt,
	})
}

func (suite *IntegrationTestSuite) sendRequest(r *http.Request, result *Response) {
	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		suite.FailNow(err.Error())
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		suite.FailNow(err.Error())
	}
}

func TestMain(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}