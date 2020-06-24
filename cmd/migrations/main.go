package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/database"
	"rampart/internal/migrations"
	"rampart/internal/secrets"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("main: migrations started")
	scr, err := secrets.NewSecrets()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Migrations.DSNParams)
	if err != nil {
		log.Fatal(err)
	}
	if err = migrations.Run(db); err != nil {
		_ = db.Close()
		log.Fatal(err)
	}
	database.CloseDatabase(db)
	log.Debug("main: migrations finished")
}
