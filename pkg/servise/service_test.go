package service

import (
	"testing"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories/mocks"
	"github.com/stretchr/testify/suite"
)

type UserTest struct {
	*models.User
	expectedError bool
	method string
}

var userSlice []*UserTest = []*UserTest{
	// SignIn test
	{
		User: &models.User{
			UserName: "Test",
			Password: "123456",
		},
		expectedError: false,
		method: "SignIn",
	},
	{
		User: &models.User{
			Email: "Test",
			Password: "123456",
		},
		expectedError: false,
		method: "SignIn",
	},
	{
		User: &models.User{
			UserName: "Test",
		},
		expectedError: true,
		method: "SignIn",
	},
	{
		User: &models.User{
			UserName: "Test",
			Password: "Wrong password",
		},
		expectedError: true,
		method: "SignIn",
	},
	// SignUp test
	{
		User: &models.User{
			UserName: "Test",
			Email: "Test",
			Password: "123456",
		},
		expectedError: false,
		method: "SignUp",
	},
	{
		User: &models.User{
			UserName: "Test",
			Email: "Test",
			Password: "12345",
		},
		expectedError: true,
		method: "SignUp",
	},

}

type serviceSuit struct {
	service *Service
	suite.Suite
	repo mocks.RepositoryUser
}

func (s *serviceSuit) SetupTest() {
	s.repo = *mocks.NewRepositoryUser(s.T())
	config, err := config.ReadConfig()
	if err != nil {
		s.FailNow(err.Error())
	}
	service := NewService(&repositories.Repository{RepositoryUser: &s.repo}, config)
	s.service = service
}

func (s *serviceSuit) TestSignUpUser() {
	for _, v := range userSlice {
		if v.method == "SignUp" {
			if v.expectedError {
				s.repo.On("SignUpUser", v.User).Return(nil)
				err := s.service.SignUpUser(v.User)
				s.Error(err)
			}
		}
	}
}

func (s *serviceSuit) TestSignInUser() {
	for _, v := range userSlice {
		if v.method == "SignIn" {
			if v.expectedError {
				s.repo.On("GetUser", v.User).Return(new(models.User), nil)
				_, err := s.service.SignInUser(v.User)
				s.Error(err)
			}
		}
	}
}

func TestStart(t *testing.T) {
	suite.Run(t, &serviceSuit{})
}
