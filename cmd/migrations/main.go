package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"rampart/internal/migrations"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("main: migrations started")
	if err := migrations.Run(); err != nil {
		log.Fatal(err)
	}
	log.Debug("main: migrations finished")
}
