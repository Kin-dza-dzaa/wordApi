package service

import (
	"testing"

	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/internal/validation"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type userTest struct {
	*models.User
	expectedError bool
	method string
}

var userSlice []*userTest = []*userTest{
	{
		User: &models.User{
			UserName: "Robot123",
			Password: "ValidPasswordIopdkfm1231231",
		},
		expectedError: false,
		method: "SignIn",
	},
	{
		User: &models.User{
			Email: "Robot123@gmail.com",
			Password: "ValidPasswordIopdkfm1231231",
		},
		expectedError: false,
		method: "SignIn",
	},
	{
		User: &models.User{
			UserName: "Robot123",
		},
		expectedError: true,
		method: "SignIn",
	},
	{
		User: &models.User{
			UserName: "##!@#$!",
			Password: "Robot123",
		},
		expectedError: true,
		method: "SignIn",
	},
	{
		User: &models.User{
		},
		expectedError: true,
		method: "SignIn",
	},
	{
		User: &models.User{
			UserName: "Test",
			Password: "1234",
		},
		expectedError: true,
		method: "SignIn",
	},
	{
		User: &models.User{
			Email: "asdlasd",
		},
		expectedError: true,
		method: "SignIn",
	},
	{
		User: &models.User{
			UserName: "Robot123",
			Email: "Robot123@gmail.com",
			Password: "ValidPasswordIopdkfm1231231",
		},
		expectedError: false,
		method: "SignUp",
	},
	{
		User: &models.User{
			UserName: "!@@$",
			Email: "Robot123@gmail.com",
			Password: "12345",
		},
		expectedError: true,
		method: "SignUp",
	},
	{
		User: &models.User{
			UserName: "Robot123",
			Email: "Robot123",
			Password: "12345",
		},
		expectedError: true,
		method: "SignUp",
	},
	{
		User: &models.User{
			UserName: "Robot123",
			Email: "Robot123@gmail.com",
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
	validator, err := validation.InitValidators()
	if err != nil {
		s.FailNow(err.Error())
	}
	service := NewService(&repositories.Repository{RepositoryUser: &s.repo}, config, validator)
	s.service = service
}

func (s *serviceSuit) TestSignUpUser() {
	for _, v := range userSlice {
		if v.method == "SignUp" {
			if v.expectedError {
				s.repo.On("SignUpUser", v.User).Return(nil)
				err := s.service.SignUpUser(v.User)
				s.Error(err)
			} else {
				s.repo.On("SignUpUser", v.User).Return(nil)
				err := s.service.SignUpUser(v.User)
				s.Nil(err)
			}
		}
	}
}

func (s *serviceSuit) TestSignInUser() {
	for _, v := range userSlice {
		if v.method == "SignIn" {
			if v.expectedError {
				_, err := s.service.SignInUser(v.User)
				s.Error(err)
			} else {
				hash, err := bcrypt.GenerateFromPassword([]byte(v.Password), 14)
				if err != nil {
					s.T().FailNow()
				}
				s.repo.On("GetUser", v.User).Return(&models.User{UserId: uuid.New(), Password: string(hash)}, nil)
				_, err = s.service.SignInUser(v.User)
				s.Nil(err)
			}
		}
	}
}

func TestUserSevice(t *testing.T) {
	suite.Run(t, &serviceSuit{})
}
