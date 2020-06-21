package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("main: mining started")
	if err := mining.Schedule(); err != nil {
		log.Fatal(err)
	}
	log.Debug("main: mining finished")
}
