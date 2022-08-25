package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserId   uuid.UUID		`json:"-"`
	UserName string    	`json:"user_name"`
	Email     string    	`json:"email"`
	Password  string    	`json:"password"`
	Time      time.Time     `json:"-"`
}