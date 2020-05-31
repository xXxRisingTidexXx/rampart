package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining/config"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("mining: started")
	mining, err := config.NewMining()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(mining)
	//prospector := domria.NewProspector(mining.Secondary)
	//if err := prospector.Prospect(); err != nil {
	//	log.Fatal(err)
	//}
	log.Debug("mining: finished")
}
