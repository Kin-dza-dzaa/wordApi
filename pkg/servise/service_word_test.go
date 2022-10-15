package service

// Right now there's no service logic to test

import (
	"os"
	"testing"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type testWordSuite struct {
	service Service
	suite.Suite
}

func (suite *testWordSuite) SetupSuite() {
	repo := mocks.NewRepository(suite.T())
	config, err := config.ReadConfig()
	if err != nil {
		suite.FailNow(err.Error())
	}
	myLogger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	suite.service = NewService(repo, config, &myLogger)
}

func TestWordService(t *testing.T) {
	suite.Run(t, new(testWordSuite))
}
