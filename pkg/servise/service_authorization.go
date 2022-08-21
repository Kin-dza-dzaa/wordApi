package service

import (
	"errors"
	"time"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	config *config.Config
	repository repository.Repository
}

func (a *service) SignUpUser(user *models.User) error {
	if len(user.Password) <= 0 {
		return errors.New("too short password")
	}
	if err := a.hashPassword(user); err != nil {
		return err
	}
	user.User_id = uuid.New()
	user.Time = time.Now()
	return a.repository.SignUpUser(user)
}

func (a *service) GenerateToken(user *models.User) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, models.MyJwtClaims{UserId: user.User_id.String()}).SignedString(a.config.JWTString)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {return a.config.JWTString, nil})
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

func (a *service) SignInUser(user *models.User) (uuid.UUID, error) {
	tempUser, err := a.repository.GetUser(user.Email)
	if err != nil {
		return uuid.UUID{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(user.Password))
	if err != nil {
		return tempUser.User_id, nil
	}
	return uuid.UUID{}, err 
}

func (a *service) hashPassword(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
