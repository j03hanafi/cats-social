package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

var ErrCatNotFound = errors.New("cat not found")

type CatRace string

const (
	Persian          CatRace = "Persian"
	MaineCoon                = "Maine Coon"
	Siamese                  = "Siamese"
	Ragdoll                  = "Ragdoll"
	Bengal                   = "Bengal"
	Sphynx                   = "Sphynx"
	BritishShorthair         = "British Shorthair"
	Abyssinian               = "Abyssinian"
	ScottishFold             = "Scottish Fold"
	Birman                   = "Birman"
)

func (r CatRace) Validate() error {
	validRace := map[CatRace]struct{}{
		Persian:          {},
		MaineCoon:        {},
		Siamese:          {},
		Ragdoll:          {},
		Bengal:           {},
		Sphynx:           {},
		BritishShorthair: {},
		Abyssinian:       {},
		ScottishFold:     {},
		Birman:           {},
	}
	defer clear(validRace)

	if _, ok := validRace[r]; !ok {
		return fmt.Errorf("invalid cat race: %s", r)
	}

	return nil
}

type CatSex string

const (
	Male   CatSex = "male"
	Female CatSex = "female"
)

func (s CatSex) Validate() error {
	validSex := map[CatSex]struct{}{
		Male:   {},
		Female: {},
	}
	defer clear(validSex)

	if _, ok := validSex[s]; !ok {
		return fmt.Errorf("invalid cat sex: %s", s)
	}

	return nil
}

type Cat struct {
	ID          ulid.ULID
	Name        string
	Race        CatRace
	Sex         CatSex
	AgeInMonth  int
	Description string
	UserID      ulid.ULID
	HasMatched  bool
	ImageUrls   []string
	CreatedAt   time.Time
}
