package migrations

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run() error {
	//"postgres://postgres:postgres@localhost:5432/rampart?sslmode=disable&connect_timeout=10"
	db, err := sql.Open("postgres", os.Getenv("RAMPART_DSN"))
	if err != nil {
		log.Fatalf("cmd: upgrade failed to connect to the db, %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("main: upgrade failed to close the db, %v", err)
		}
		log.Debug("main: migrations finished")
	}()
	if err = db.Ping(); err != nil {
		log.Fatalf("cmd: upgrade failed to ping the db, %v", err)
	}
	return nil
}
