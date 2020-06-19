package migrations

import (
	"database/sql"
	"fmt"
)

// TODO: postgres://postgres:postgres@localhost:5432/rampart
func Run() error {
	migrations, err := newMigrations()
	if err != nil {
		return err
	}
	dsn, err := getDSN(migrations.QueryParams)
	if err != nil {
		return err
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("migrations: failed to open the db, %v", err)
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("migrations: failed to ping the db, %v", err)
	}
	versions, err := listVersions()
	if err != nil {
		_ = db.Close()
		return err
	}
	migrator, err := newMigrator(db)
	if err != nil {
		_ = db.Close()
		return err
	}
	latest, err := migrator.ensureVersions()
	if err != nil {
		_ = db.Close()
		return err
	}
	for _, version := range versions {
		if version.id > latest {
			if err = migrator.applyVersion(version); err != nil {
				_ = db.Close()
				return err
			}
		}
	}
	if err = migrator.commit(); err != nil {
		_ = db.Close()
		return err
	}
	if err = db.Close(); err != nil {
		return fmt.Errorf("migrations: failed to close the db, %v", err)
	}
	return nil
}
