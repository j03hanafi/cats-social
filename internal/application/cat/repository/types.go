package repository

import (
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
)

type cat struct {
	ID          ulid.ULID
	Name        string
	Race        string
	Sex         string
	AgeInMonth  int
	Description string
	UserID      ulid.ULID
	HasMatched  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

type catImages struct {
	ID        ulid.ULID
	ImageURL  string
	CatID     ulid.ULID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

func (i catImages) tableName() string {
	return "cat_images"
}
