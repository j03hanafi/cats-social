package repository

import (
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
)

type user struct {
	ID        ulid.ULID
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
