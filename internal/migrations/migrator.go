package migrations

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// TODO: inject tx instead of db
func newMigrator(db *sql.DB) (*migrator, error) {
	return &migrator{db}, nil
}

type migrator struct {
	db *sql.DB
}

// TODO: add latest reading
func (migrator *migrator) ensureVersions() (int64, error) {
	start := time.Now()
	tx, err := migrator.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("migrations: migrator failed to init the versions, %v", err)
	}
	if _, err = tx.Exec("create table if not exists versions(id bigint not null)"); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("migrations: migrator failed to create the versions, %v", err)
	}
	count := 0
	if err = tx.QueryRow("select count(*) from versions").Scan(&count); err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("migrations: migrator failed to read the versions, %v", err)
	}
	if 1 < count {
		_ = tx.Rollback()
		return 0, fmt.Errorf("migrations: migrator got multiple rows in the versions, %d", count)
	}
	if count == 0 {
		if _, err = tx.Exec("insert into versions values (0)"); err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("migrations: migrator failed to set the zero version, %v", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("migrations: migrator failed to commit the versions, %v", err)
	}
	log.Debugf("migrations: migrator ensured the versions (%.3fs)", time.Since(start).Seconds())
	return 0, nil
}

func (migrator *migrator) applyVersion(version *version) error {
	return nil
}

func (migrator *migrator) commit() error {
	return nil
}
