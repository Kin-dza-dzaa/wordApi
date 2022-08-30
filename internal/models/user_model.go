package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserId   uuid.UUID		`json:"-"`
	UserName string    		`json:"user_name,omitempty"`
	Email     string    	`json:"email,omitempty"`
	Password  string    	`json:"password"`
	Time      time.Time     `json:"-"`
}