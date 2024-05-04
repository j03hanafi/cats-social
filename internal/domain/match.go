package domain

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ErrCatGenderNotMatch = errors.New("you can't match with the same sex")
	ErrCatAlreadyMatched = errors.New("cat already matched")
	ErrCatSameOwner      = errors.New("you can't match with your own cat")
	ErrMatchNotFound     = errors.New("match not found")
	ErrMatchNotValid     = errors.New("match not valid")
)

type Match struct {
	ID         ulid.ULID
	MatchCatID ulid.ULID
	UserCatID  ulid.ULID
	Message    string
	CreatedAt  time.Time
	DeletedAt  time.Time
}

type DetailMatch struct {
	Match
	Issuer   User
	Receiver User
	MatchCat Cat
	UserCat  Cat
}
