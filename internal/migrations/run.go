package migrations

import (
	"database/sql"
	"fmt"
	"os"
)

func Run() (err error) {
	//"postgres://postgres:postgres@localhost:5432/rampart?sslmode=disable&connect_timeout=10"
	dsn := os.Getenv("RAMPART_DSN")
	if dsn == "" {
		return fmt.Errorf("migrations: dsn isn't configured")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("migrations: failed to connect to the db, %v", err)
	}
	defer func() {
		if closingErr := db.Close(); closingErr != nil && err == nil {
			err = fmt.Errorf("migrations: failed to close the db, %v", closingErr)
		}
	}()
	if err = db.Ping(); err != nil {
		return fmt.Errorf("migrations: failed to ping the db, %v", err)
	}

	return nil
}
