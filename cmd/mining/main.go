package main

import (
	"flag"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/database"
	"rampart/internal/mining"
	"rampart/internal/mining/metrics"
	"rampart/internal/secrets"
)

func main() {
	isOnce := flag.Bool("once", false, "Execute a single workflow instead of the whole schedule")
	alias := flag.String("miner", "", "Desired miner alias")
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.Debug("main: mining started")
	scr, err := secrets.NewSecrets()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Mining.DSNParams)
	if err != nil {
		log.Fatal(err)
	}
	miner, err := mining.FindMiner(*alias, cfg.Mining.Miners, db)
	if err != nil {
		_ = db.Close()
		log.Fatal(err)
	}
	if *isOnce {
		miner.Run()
	} else {
		scheduler := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
		if _, err = scheduler.AddJob(miner.Spec(), miner); err != nil {
			_ = db.Close()
			log.Fatalf("main: mining failed to schedule, %v", err)
		}
		metrics.RunServer(cfg.Mining.Metrics.Server)
		scheduler.Run()
	}
	database.CloseDatabase(db)
	log.Debug("main: mining finished")
}
