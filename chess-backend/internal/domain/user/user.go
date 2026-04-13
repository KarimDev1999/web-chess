package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string
	Email     string
	Password  string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, hashedPassword, username string) *User {
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  hashedPassword,
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
