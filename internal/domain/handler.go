package domain

import (
	"errors"
	"strconv"

	"github.com/oklog/ulid/v2"
	"go.uber.org/multierr"
)

const (
	InvalidRequestBodyMessage  = "Invalid Request Body"
	InternalServerErrorMessage = "Internal Server Error"
	NotFoundErrorMessage       = "Not Found"
)

type boolQueryParam string

const (
	TrueBool  boolQueryParam = "true"
	FalseBool boolQueryParam = "false"
)

func (b boolQueryParam) validate() error {
	if b != TrueBool && b != FalseBool {
		return errors.New("invalid boolean query param")
	}

	return nil
}

type QueryParam struct {
	ID         ulid.ULID      `query:"id"`
	Limit      int            `query:"limit"`
	Offset     int            `query:"offset"`
	Race       CatRace        `query:"race"`
	Sex        CatSex         `query:"sex"`
	HasMatched boolQueryParam `query:"hasMatched"`
	AgeInMonth string         `query:"ageInMonth"`
	Owned      boolQueryParam `query:"owned"`
	Search     string         `query:"search"`
}

func (p *QueryParam) Validate() error {
	var errs error

	if p.Limit == 0 {
		p.Limit = 5
	}

	if p.Race != "" {
		if err := p.Race.Validate(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if p.Sex != "" {
		if err := p.Sex.Validate(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if p.HasMatched != "" {
		if err := p.HasMatched.validate(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if err := p.validateAgeInMonth(); err != nil {
		errs = multierr.Append(errs, err)
	}

	if p.Owned != "" {
		if err := p.Owned.validate(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

func (p *QueryParam) validateAgeInMonth() error {
	if p.AgeInMonth == "" {
		return nil
	}

	// Check if the whole string is numeric (case: "%d")
	if _, err := strconv.Atoi(p.AgeInMonth); err == nil {
		p.AgeInMonth = "=" + p.AgeInMonth
		return nil
	}

	// Check for '>' or '<' followed by a numeric string (case: ">%d" or "<%d")
	if len(p.AgeInMonth) > 1 && p.AgeInMonth[0] == '>' || p.AgeInMonth[0] == '<' {
		if _, err := strconv.Atoi(p.AgeInMonth[1:]); err == nil {
			return nil
		}
	}

	return errors.New("invalid age in month query param")
}
