package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	User_id   uuid.UUID		`json:"-"`
	User_name string    	`json:"user_name"`
	Email     string    	`json:"email"`
	Password  string    	`json:"password"`
	Time      time.Time     `json:"-"`
}