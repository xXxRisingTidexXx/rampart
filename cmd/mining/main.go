package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
	"rampart/internal/mining/config"
	"rampart/internal/mining/domria"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("main: mining started")
	cfg, err := config.NewMining()
	if err != nil {
		log.Fatal(err)
	}
	prospector := domria.NewProspector(mining.Secondary, cfg.Prospectors.Domria)
	if err = prospector.Prospect(); err != nil {
		log.Fatal(err)
	}
	log.Debug("main: mining finished")
}
