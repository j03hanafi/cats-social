package domain

import (
	"errors"

	"github.com/oklog/ulid/v2"
)

var (
	DuplicateEmailError = errors.New("email already exists")
	UserNotFoundError   = errors.New("user not found")
	InvalidPassword     = errors.New("invalid password")
)

type User struct {
	ID       ulid.ULID
	Email    string
	Name     string
	Password string
}
