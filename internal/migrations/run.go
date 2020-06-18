package migrations

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run() error {
	//"postgres://postgres:postgres@localhost:5432/rampart?sslmode=disable&connect_timeout=10"
	migrations, err := newMigrations()
	if err != nil {
		return err
	}
	log.Debug(migrations)
	dsn := os.Getenv("RAMPART_DSN")
	if dsn == "" {
		return fmt.Errorf("migrations: dsn isn't configured")
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
