package domria

import (
	"database/sql"
)

func newStorer(db *sql.DB) *storer {
	return &storer{db}
}

type storer struct {
	db *sql.DB
}

func (storer *storer) storeFlats(flats []*flat) error {

	return nil
}
