package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/servise/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type testStruct struct {
	body string
	expectedError bool
	message string
}

type response struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

type TestSuite struct {
	suite.Suite
	serviceMocks *mocks.Service
	handlers *Handlers
	userId uuid.UUID
}

func (suite *TestSuite) SetupSuite() {
	suite.serviceMocks = mocks.NewService(suite.T())
	suite.handlers = NewHandlers(suite.serviceMocks)
	suite.userId = uuid.New()
} 

func (suite *TestSuite) TestGetWords() {
	var testSlice []testStruct = []testStruct{
		{
			expectedError: false,
			message: "words were sent",
		},
		{
			expectedError: true,
			message: "invalid token",
		},
	}
	for _, v := range testSlice {
		var result response
		if v.expectedError {
			r := httptest.NewRequest("GET", "/words", nil)
			w := httptest.NewRecorder()
			suite.handlers.GetWordsHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		} else {
			r := httptest.NewRequest("GET", "/words", nil)
			r = suite.addTocontext(r)
			w := httptest.NewRecorder()
			suite.serviceMocks.On("GetWords", suite.userId.String()).Return(&models.WordsGet{}, nil).Once()
			suite.handlers.GetWordsHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		}
	}	
}


func (suite *TestSuite) TestAddWordsHandler() {
	var testSlice []testStruct = []testStruct{
		{
			expectedError: false,
			message: "words were added",
			body: `{"words": [{"word": "asd", "state": 1, "collection_name": "asd"}]}`,
		},
		{
			expectedError: true,
			message: `expected ojbect: {"words": [{"word": string, "state": int, "collection_name": string}, ...]}`,
			body: "{}",
		},
	}
	for _, v := range testSlice {
		var result response
		if v.expectedError {
			r := httptest.NewRequest("POST", "/words", bytes.NewReader([]byte(v.body)))
			w := httptest.NewRecorder()
			suite.handlers.AddWordsHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		} else {
			r := httptest.NewRequest("GET", "/words", bytes.NewReader([]byte(v.body)))
			r = suite.addTocontext(r)
			w := httptest.NewRecorder()
			suite.serviceMocks.On("AddWords", mock.Anything, mock.Anything).Return([]string{}, nil).Once()
			suite.handlers.AddWordsHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *TestSuite) TestUpdateWordsHandler() {
	var testSlice []testStruct = []testStruct{
		{
			expectedError: true,
			message: "invalid token",
			body: `{}`,
		},
		{
			expectedError: true,
			message: `expected object: {"old_word": string, "new_word": string, "collection_name": string}`,
			body: `{"old_word": 1}`,
		},
		{
			expectedError: true,
			message: "invalid token",
			body: `{"old_word": "", "new_word": "", "collection_name": ""}`,
		},
		{
			expectedError: false,
			message: " was updated to ",
			body: `{"old_word": "", "new_word": "", "collection_name": ""}`,
		},
	}
	for _, v := range testSlice {
		var result response
		if v.expectedError {
			r := httptest.NewRequest("POST", "/words", bytes.NewReader([]byte(v.body)))
			w := httptest.NewRecorder()
			suite.handlers.UpdateWordHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		} else {
			r := httptest.NewRequest("GET", "/words", bytes.NewReader([]byte(v.body)))
			r = suite.addTocontext(r)
			w := httptest.NewRecorder()
			suite.serviceMocks.On("UpdateWord", mock.Anything, mock.Anything).Return(nil).Once()
			suite.handlers.UpdateWordHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *TestSuite) TestDeleteWordHandler() {
	var testSlice []testStruct = []testStruct{
		{
			expectedError: true,
			message: "invalid token",
			body: `{"words": []}`,
		},
		{
			expectedError: true,
			message: `expected object: {"words": [{"collection_name": string, "word": string}...]}`,
			body: `{"old_word": 1}`,
		},
		{
			expectedError: true,
			message: "invalid token",
			body: `{"words": []}`,
		},
	}
	for _, v := range testSlice {
		var result response
		if v.expectedError {
			r := httptest.NewRequest("POST", "/words", bytes.NewReader([]byte(v.body)))
			w := httptest.NewRecorder()
			suite.handlers.DeleteWordHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		} else {
			r := httptest.NewRequest("GET", "/words", bytes.NewReader([]byte(v.body)))
			r = suite.addTocontext(r)
			w := httptest.NewRecorder()
			suite.serviceMocks.On("DeleteWords", mock.Anything, mock.Anything).Once()
			suite.handlers.DeleteWordHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *TestSuite) TestUpdateStateHandler() {
	var testSlice []testStruct = []testStruct{
		{
			expectedError: true,
			message: "invalid token",
			body: `{"words": []}`,
		},
		{
			expectedError: true,
			message: `expected object: {"collection_name": string, "words": object[{"word": string, "new_state": number}]}`,
			body: `{}`,
		},
	}
	for _, v := range testSlice {
		var result response
		if v.expectedError {
			r := httptest.NewRequest("POST", "/words", bytes.NewReader([]byte(v.body)))
			w := httptest.NewRecorder()
			suite.handlers.UpdateStateHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		} else {
			r := httptest.NewRequest("GET", "/words", bytes.NewReader([]byte(v.body)))
			r = suite.addTocontext(r)
			w := httptest.NewRecorder()
			suite.serviceMocks.On("UpdateState", mock.Anything, mock.Anything).Once()
			suite.handlers.UpdateStateHandler().ServeHTTP(w, r)
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				suite.FailNow(err.Error())
			}
			suite.Equal(v.message, result.Message)
		}
	}	
}

func (suite *TestSuite) addTocontext(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, KeyForToken("user_id"), suite.userId.String())
	return r.WithContext(ctx)
}

func TestRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}