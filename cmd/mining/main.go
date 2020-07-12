package main

import (
	"flag"
	_ "github.com/lib/pq"
	gocron "github.com/robfig/cron/v3"
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
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("miner", *alias)
	scr, err := secrets.NewSecrets()
	if err != nil {
		entry.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Mining.DSNParams)
	if err != nil {
		entry.Fatal(err)
	}
	gatherer := metrics.NewGatherer(*alias, db)
	miner, err := mining.FindMiner(*alias, cfg.Mining.Miners, db, gatherer, entry)
	if err != nil {
		_ = db.Close()
		entry.Fatal(err)
	}
	if *isOnce {
		miner.Run()
	} else {
		cron := gocron.New()
		if _, err = cron.AddJob(miner.Spec(), miner); err != nil {
			_ = db.Close()
			entry.Fatalf("main: mining failed to run, %v", err)
		}
		metrics.RunServer(miner.Port(), cfg.Mining.Server, entry)
		cron.Run()
	}
	if err = database.CloseDatabase(db); err != nil {
		entry.Fatal(err)
	}
}
