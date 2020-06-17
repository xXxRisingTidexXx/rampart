package migrations

import (
	"database/sql"
	"fmt"
	"os"
)

func Run() error {
	//"postgres://postgres:postgres@localhost:5432/rampart?sslmode=disable&connect_timeout=10"
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
	migrator := newMigrator(db)
	if err = migrator.ensureVersions(); err != nil {
		_ = db.Close()
		return err
	}

	if err = db.Close(); err != nil {
		return fmt.Errorf("migrations: failed to close the db, %v", err)
	}
	return nil
}
