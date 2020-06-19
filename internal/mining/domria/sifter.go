package domria

import (
	"database/sql"
)

func newSifter(db *sql.DB) *sifter {
	return &sifter{db}
}

type sifter struct {
	db *sql.DB
}

func (sifter *sifter) siftFlats(flats []*flat) []*flat {
	return flats
}
