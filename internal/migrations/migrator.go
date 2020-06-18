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
		if _, err := migrator.tx.Exec("insert into versions values ($1)", id); err != nil {
			_ = migrator.tx.Rollback()
			return 0, fmt.Errorf("migrations: migrator failed to create version %d, %v", id, err)
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
	start := time.Now()
	query, err := version.load()
	if err != nil {
		_ = migrator.tx.Rollback()
		return err
	}
	if _, err = migrator.tx.Exec(query); err != nil {
		_ = migrator.tx.Rollback()
		return fmt.Errorf("migrations: migrator failed to apply version %d, %v", version.id, err)
	}
	if _, err := migrator.tx.Exec("update versions set id = $1", version.id); err != nil {
		_ = migrator.tx.Rollback()
		return fmt.Errorf("migrations: migrator failed to update version %d, %v", version.id, err)
	}
	log.Debugf("migrations: migrator applied version %d (%.3fs)", version.id, time.Since(start).Seconds())
	return nil
}

func (migrator *migrator) commit() error {
	start := time.Now()
	if err := migrator.tx.Commit(); err != nil {
		return fmt.Errorf("migrations: migrator failed to commit the changes, %v", err)
	}
	log.Debugf("migrations: migrator committed (%.3fs)", time.Since(start).Seconds())
	return nil
}
