package repository

import (
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
)

type match struct {
	ID         ulid.ULID
	MatchCatID ulid.ULID
	UserCatID  ulid.ULID
	Message    string
	IssuerID   ulid.ULID
	ReceiverID ulid.ULID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime
}
