package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserId    		uuid.UUID	`json:"-"`
	UserName  		string    	`json:"user_name,omitempty" validate:"omitempty,userval"`
	Email     		string    	`json:"email,omitempty" validate:"omitempty,email"`
	Password  		string    	`json:"password" validate:"required,passval"`
	Time      		time.Time   `json:"-"`
}
