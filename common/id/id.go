package id

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func New() ulid.ULID {
	timestamp := ulid.Timestamp(time.Now())
	newID, err := ulid.New(timestamp, rand.Reader)
	if err != nil {
		return ulid.ULID{}
	}
	return newID
}
