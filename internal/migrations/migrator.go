package migrations

import (
	"database/sql"
)

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db}
}

type Migrator struct {
	db *sql.DB
}

func (migrator *Migrator) EnsureVersions() error {
	return nil
}
