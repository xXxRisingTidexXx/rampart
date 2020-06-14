package migrations

import (
	"database/sql"
)

func NewUpgrader(db *sql.DB) *Upgrader {
	return &Upgrader{db}
}

type Upgrader struct {
	db *sql.DB
}

func (upgrader *Upgrader) EnsureVersions() error {
	return nil
}
