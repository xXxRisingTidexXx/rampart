package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func main() {
	_ = flag.String(
		"miner",
		"",
		"Set a concrete miner name to run it once; leave the field blank to up the whole messis",
	)
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "messis")
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := sql.Open("postgres", c.Messis.DSN)
	if err != nil {
		entry.Fatalf("main: messis failed to open the db, %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		entry.Fatalf("main: messis failed to ping the db, %v", err)
	}


	if err = db.Close(); err != nil {
		entry.Fatalf("main: messis failed to close the db, %v", err)
	}
}
