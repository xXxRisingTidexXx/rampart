package main

import (
	"flag"
	_ "github.com/lib/pq"
	gocron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/database"
	"github.com/xXxRisingTidexXx/rampart/internal/mining"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/secrets"
)

func main() {
	isOnce := flag.Bool("once", false, "Execute a single workflow instead of the whole schedule")
	alias := flag.String("miner", "", "Desired miner alias")
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	logger := logging.NewLogger(*alias)
	scr, err := secrets.NewSecrets()
	if err != nil {
		logger.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Mining.DSNParams)
	if err != nil {
		logger.Fatal(err)
	}
	gatherer := metrics.NewGatherer(*alias, db)
	miner, err := mining.FindMiner(*alias, cfg.Mining.Miners, db, gatherer, logger)
	if err != nil {
		_ = db.Close()
		logger.Fatal(err)
	}
	if *isOnce {
		miner.Run()
	} else {
		cron := gocron.New()
		if _, err = cron.AddJob(miner.Spec(), miner); err != nil {
			_ = db.Close()
			logger.Fatalf("main: mining failed to run, %v", err)
		}
		metrics.RunServer(miner.Port(), cfg.Mining.Server, logger)
		cron.Run()
	}
	if err = database.CloseDatabase(db); err != nil {
		logger.Fatal(err)
	}
}
