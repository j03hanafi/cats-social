package domain

import "github.com/oklog/ulid/v2"

type User struct {
	ID       ulid.ULID
	Email    string
	Name     string
	Password string
}
