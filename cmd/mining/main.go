package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
	"rampart/pkg/mining/domria"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("mining: started")
	prospector := domria.NewProspector(mining.Secondary)
	if err := prospector.Prospect(); err != nil {
		log.Fatal(err)
	}
	log.Debug("mining: finished")
}
