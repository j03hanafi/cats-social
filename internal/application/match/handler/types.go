package handler

import (
	"errors"

	"github.com/oklog/ulid/v2"
	"go.uber.org/multierr"

	"cats-social/internal/domain"
)

const (
	successMatchMessage        = "Match created successfully"
	successGetMatchMessage     = "Success"
	successApproveMatchMessage = "Match approved successfully"
	successRejectMatchMessage  = "Match rejected successfully"
	successDeleteMatchMessage  = "Match deleted successfully"
)

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type matchRequest struct {
	MatchCatId ulid.ULID `json:"matchCatId"`
	UserCatId  ulid.ULID `json:"userCatId"`
	Message    string    `json:"message"`
}

func (a matchRequest) validate() error {
	var errs error

	if a.MatchCatId.String() == "" {
		errs = multierr.Append(errs, errors.New("matchCatId is required"))
	}
	if a.UserCatId.String() == "" {
		errs = multierr.Append(errs, errors.New("userCatId is required"))
	}

	if a.Message != "" && len(a.Message) < 5 || len(a.Message) > 120 {
		errs = multierr.Append(errs, errors.New("message length must be between 5 and 120"))
	}
	if a.Message == "" {
		errs = multierr.Append(errs, errors.New("message is required"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

type detailMatchResponse struct {
	ID             ulid.ULID `json:"id"`
	IssuedBy       issuedBy  `json:"issuedBy"`
	MatchCatDetail catDetail `json:"matchCatDetail"`
	UserCatDetail  catDetail `json:"userCatDetail"`
	Message        string    `json:"message"`
	CreatedAt      string    `json:"createdAt"`
}

type issuedBy struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type catDetail struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Race        domain.CatRace `json:"race"`
	Sex         domain.CatSex  `json:"sex"`
	AgeInMonth  int            `json:"ageInMonth"`
	Description string         `json:"description"`
	ImageUrls   []string       `json:"imageUrls"`
	HasMatched  bool           `json:"hasMatched"`
	CreatedAt   string         `json:"createdAt"`
}

type approvalMatchRequest struct {
	MatchID ulid.ULID `json:"matchId"`
}

func (a approvalMatchRequest) validate() error {
	var emptyULID ulid.ULID
	if a.MatchID == emptyULID {
		return errors.New("matchId is required")
	}
	return nil
}
