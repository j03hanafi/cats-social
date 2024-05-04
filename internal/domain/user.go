package domain

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	UserFromToken = "loggedInUser"
)

var (
	DuplicateEmailError = errors.New("email already exists")
	UserNotFoundError   = errors.New("user not found")
	InvalidPassword     = errors.New("invalid password")
)

type User struct {
	ID        ulid.ULID
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
}
