package service

// there're no service logic in that layer right now

import (
	"testing"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories/mocks"
	"github.com/stretchr/testify/suite"
	"github.com/Kin-dza-dzaa/wordApi/internal/validation"
)

type testWordSuite struct {
	service *Service
	suite.Suite
}

func (w *testWordSuite) SetupSuite() {
	repo := mocks.NewRepositoryUser(w.T())
	config, err := config.ReadConfig()
	if err != nil {
		w.FailNow(err.Error())
	}
	validator, err := validation.InitValidators()
	if err != nil {
		w.FailNow(err.Error())
	}
	w.service = NewService(&repositories.Repository{RepositoryUser: repo}, config, validator)
}

func TestWordService(t *testing.T) {
	suite.Run(t, new(testWordSuite))
}
