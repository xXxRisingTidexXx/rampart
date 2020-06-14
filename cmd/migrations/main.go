package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/migrations/cmd"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("main: migrations started")
	cmd.Execute()
	log.Debug("main: migrations finished")
}
