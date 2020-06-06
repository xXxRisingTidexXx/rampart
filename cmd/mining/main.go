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
	exitIfError(err)
	prospector := domria.NewProspector(mining.Secondary, config.Prospectors.Domria)
	exitIfError(prospector.Prospect())
	log.Debug("main: mining finished")
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
