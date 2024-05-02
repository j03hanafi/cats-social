package handler

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"go.uber.org/multierr"
)

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r registerRequest) validate() error {
	var errs error

	if r.Email != "" && !govalidator.IsEmail(r.Email) {
		errs = multierr.Append(errs, errors.New("invalid email format"))
	}
	if r.Email == "" {
		errs = multierr.Append(errs, errors.New("email is required"))
	}

	if r.Name != "" && len(r.Name) < 5 || len(r.Name) > 50 {
		errs = multierr.Append(errs, errors.New("name must be between 5 and 50 characters"))
	}
	if r.Name == "" {
		errs = multierr.Append(errs, errors.New("name is required"))
	}

	if r.Password != "" && len(r.Password) < 5 || len(r.Password) > 15 {
		errs = multierr.Append(errs, errors.New("password must be between 5 and 15 characters"))
	}
	if r.Password == "" {
		errs = multierr.Append(errs, errors.New("password is required"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var invalidRequestBody = baseResponse{
	Message: "Invalid Request Body",
}
