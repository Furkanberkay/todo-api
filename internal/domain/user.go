package domain

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, auth *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	gorm.Model
	Name           string `gorm:"not null;size:100"`
	Surname        string `gorm:"not null;size:100"`
	Username       string `gorm:"not null;unique;size:100"`
	Email          string `gorm:"not null;unique;size:100"`
	HashedPassword string `gorm:"not null;size:300"`
	Todos          []Todo `gorm:"foreignKey:UserID"`
}

var ErrUserAlreadyExists = errors.New("user already exist")
var InternalError = errors.New("internal Error")
var InvalidInput = errors.New("invalid input")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
