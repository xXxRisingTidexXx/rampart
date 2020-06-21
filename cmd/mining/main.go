package main

import (
	"flag"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
)

func main() {
	isOnce := flag.Bool("once", false, "Execute a single workflow instead of the whole schedule")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.Debug("main: mining started")
	function := mining.Schedule
	if *isOnce {
		function = mining.Run
	}
	if err := function(); err != nil {
		log.Fatal(err)
	}
	log.Debug("main: mining finished")
}
