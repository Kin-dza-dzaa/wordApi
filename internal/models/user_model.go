package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserId           uuid.UUID `json:"-"`
	CsrfToken        string    `json:"-"`
	Jwt              string    `json:"-"`
}
