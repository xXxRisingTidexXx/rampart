package telegram

import (
	"database/sql"
)

type Observer struct {
	db *sql.DB
}

func (observer *Observer) ObserveLookups() ([]Lookup, error) {
	return nil, nil
}
