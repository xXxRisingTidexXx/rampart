package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/database"
	"github.com/xXxRisingTidexXx/rampart/internal/migrations"
	"github.com/xXxRisingTidexXx/rampart/internal/secrets"
)

// TODO: try to use this instead: https://github.com/golang-migrate/migrate .
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "migrations")
	scr, err := secrets.NewSecrets()
	if err != nil {
		entry.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Migrations.DSNParams)
	if err != nil {
		entry.Fatal(err)
	}
	if err := migrations.Run(db); err != nil {
		_ = db.Close()
		entry.Fatal(err)
	}
	if err := database.CloseDatabase(db); err != nil {
		entry.Fatal(err)
	}
}
