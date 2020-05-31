package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
	"rampart/pkg/mining/domria"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("prospector started")
	prospector := domria.NewProspector()
	if err := prospector.Prospect(mining.Secondary); err != nil {
		log.Fatal(err)
	}
	log.Info("prospector finished")
}
