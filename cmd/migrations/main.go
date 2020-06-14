package main

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/migrations/cmd"
)

func main() {
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
}
