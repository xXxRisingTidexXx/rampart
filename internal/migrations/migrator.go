package migrations

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func newMigrator(db *sql.DB) *migrator {
	return &migrator{db}
}

type migrator struct {
	db *sql.DB
}

func (migrator *migrator) ensureVersions() error {
	log.Debug("migrations: migrator started versions transaction")
	tx, err := migrator.db.Begin()
	if err != nil {
		return fmt.Errorf("migrations: migrator failed to init the versions transaction, %v", err)
	}
	_, err = tx.Exec("create table if not exists versions(version bigint not null default 0)")
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("migrations: migrator failed to create versions, %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("migrations: migrator failed to commit the versions transaction, %v", err)
	}
	log.Debug("migrations: migrator finished versions transaction")
	return nil
}
