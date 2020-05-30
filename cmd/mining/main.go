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
	prospector.Prospect("Київська", "Київ", mining.Primary)
	log.Info("prospector finished")
}
