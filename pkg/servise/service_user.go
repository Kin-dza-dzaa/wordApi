package service

import (
	"errors"
	"time"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
)

type serviceUser struct {
	config     *config.Config
	repository *repositories.Repository
	validator *validator.Validate
}

func (a *serviceUser) SignUpUser(user *models.User) error {
	if err := a.validator.Struct(user); err != nil {
		return err
	}
	if err := a.hashPassword(user); err != nil {
		return err
	}
	user.UserId = uuid.New()
	user.Time = time.Now()
	if a.repository.SignUpUser(user) != nil {
		return errors.New("user already exists")
	}
	return nil
}

func (a *serviceUser) SignInUser(user *models.User) (string, error) {
	if err := a.validator.Struct(user); err != nil {
		return "", err
	}
	dbUser, err := a.repository.GetUser(user)
	if err != nil {
		return "", errors.New("user isn't registrated")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("wrong password")
	}
	token, err := a.generateToken(dbUser)
	if  err != nil {
		return "", errors.New("internal server error")
	}
	return token, nil
}

func (a *serviceUser) generateToken(user *models.User) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.MyJwtClaims{UserId: user.UserId.String()}).SignedString([]byte(a.config.JWTString))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *serviceUser) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {return []byte(a.config.JWTString), nil})
	if err != nil {
		return "", err
	}
	mapClaims, okClaims := token.Claims.(jwt.MapClaims)
	user_id, okMap := mapClaims["user_id"].(string)
	if okClaims && okMap {
		return user_id, nil
	}
	return "", errors.New("user isn't verified by server")
}

func (a *serviceUser) hashPassword(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
