package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
	"rampart/pkg/mining/configs"
	"rampart/pkg/mining/domria"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("mining: started")
	config, err := configs.NewMining()
	if err != nil {
		log.Fatal(err)
	}
	prospector := domria.NewProspector(mining.Secondary, config.UserAgent, config.Prospectors.Domria)
	if err := prospector.Prospect(); err != nil {
		log.Fatal(err)
	}
	log.Debug("mining: finished")
}
