package migrations

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func newMigrator(db *sql.DB) (*migrator, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("migrations: failed to acquire a transaction, %v", err)
	}
	return &migrator{tx}, nil
}

type migrator struct {
	tx *sql.Tx
}

func (migrator *migrator) ensureVersions() (int64, error) {
	start := time.Now()
	_, err := migrator.tx.Exec("create table if not exists versions(id bigint not null)")
	if err != nil {
		_ = migrator.tx.Rollback()
		return 0, fmt.Errorf("migrations: migrator failed to create the versions, %v", err)
	}
	count := 0
	if err = migrator.tx.QueryRow("select count(*) from versions").Scan(&count); err != nil {
		_ = migrator.tx.Rollback()
		return 0, fmt.Errorf("migrations: migrator failed to read the versions, %v", err)
	}
	if 1 < count {
		_ = migrator.tx.Rollback()
		return 0, fmt.Errorf("migrations: migrator got multiple rows in the versions, %d", count)
	}
	id := int64(0)
	if count == 0 {
		if _, err = migrator.tx.Exec("insert into versions values ($1)", id); err != nil {
			_ = migrator.tx.Rollback()
			return 0, fmt.Errorf("migrations: migrator failed to insert the zero version, %v", err)
		}
	} else {
		if err = migrator.tx.QueryRow("select id from versions").Scan(&id); err != nil {
			_ = migrator.tx.Rollback()
			return 0, fmt.Errorf("migrations: migrator failed to read the latest version, %v", err)
		}
	}
	log.Debugf("migrations: migrator ensured the versions (%.3fs)", time.Since(start).Seconds())
	return id, nil
}

func (migrator *migrator) applyVersion(version *version) error {
	return nil
}

func (migrator *migrator) commit() error {
	if err := migrator.tx.Commit(); err != nil {
		return fmt.Errorf("migrations: migrator failed to commit the changes, %v", err)
	}
	return nil
}
