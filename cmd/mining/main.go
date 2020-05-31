package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
	"rampart/pkg/mining/domria"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("prospector started")
	prospector := domria.NewProspector()
	if err := prospector.Prospect(mining.Secondary); err != nil {
		log.Fatal(err)
	}
	log.Debug("prospector finished")
}
