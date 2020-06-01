package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
	"rampart/internal/mining/domria"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("main: mining started")
	config, err := configs.NewMining()
	if err != nil {
		log.Fatal(err)
	}
	prospector := domria.NewProspector(mining.Secondary, config.UserAgent, config.Prospectors.Domria)
	if err := prospector.Prospect(); err != nil {
		log.Fatal(err)
	}
	log.Debug("main: mining finished")
}
