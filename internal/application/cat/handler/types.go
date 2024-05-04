package handler

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"go.uber.org/multierr"

	"cats-social/internal/domain"
)

const (
	successAddCatMessage    = "Cat added successfully"
	successListCatMessage   = "Success"
	successUpdateCatMessage = "Cat updated successfully"
	successDeleteCatMessage = "Cat deleted successfully"
)

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type imageUrls []string

func (i imageUrls) validate() error {
	if len(i) < 1 {
		return errors.New("imageUrls is required")
	}

	var errs error
	for idx, img := range i {
		if img != "" && !govalidator.IsURL(img) {
			errs = multierr.Append(errs, fmt.Errorf("invalid image URL [%d]: %s", idx, img))
		}
		if img == "" {
			errs = multierr.Append(errs, fmt.Errorf("image URL [%d] is required", idx))
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

type catRequest struct {
	Name        string         `json:"name"`
	Race        domain.CatRace `json:"race"`
	Sex         domain.CatSex  `json:"sex"`
	AgeInMonth  int            `json:"ageInMonth"`
	Description string         `json:"description"`
	ImageUrls   imageUrls      `json:"imageUrls"`
}

func (a catRequest) validate() error {
	var errs error

	if a.Name != "" && len(a.Name) < 1 || len(a.Name) > 30 {
		errs = multierr.Append(errs, errors.New("name must be between 1 and 30 characters"))
	}
	if a.Name == "" {
		errs = multierr.Append(errs, errors.New("name is required"))
	}

	if err := a.Race.Validate(); err != nil {
		errs = multierr.Append(errs, err)
	}

	if err := a.Sex.Validate(); err != nil {
		errs = multierr.Append(errs, err)
	}

	if a.AgeInMonth < 1 || a.AgeInMonth > 120082 {
		errs = multierr.Append(errs, errors.New("age must be between 1 and 120082 months"))
	}

	if a.Description != "" && len(a.Description) < 1 || len(a.Description) > 200 {
		errs = multierr.Append(errs, errors.New("description must be between 1 and 200 characters"))
	}
	if a.Description == "" {
		errs = multierr.Append(errs, errors.New("description is required"))
	}

	if err := a.ImageUrls.validate(); err != nil {
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		return errs
	}

	return nil
}

type addCatResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type listCatResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Race        domain.CatRace `json:"race"`
	Sex         domain.CatSex  `json:"sex"`
	AgeInMonth  int            `json:"ageInMonth"`
	Description string         `json:"description"`
	ImageUrls   imageUrls      `json:"imageUrls"`
	HasMatched  bool           `json:"hasMatched"`
	CreatedAt   string         `json:"createdAt"`
}
