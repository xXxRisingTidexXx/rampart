package main

import (
	"flag"
	gocron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/database"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/secrets"
)

// TODO: set service label for various logs.
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
	jobs := map[string]gocron.Job{
		cfg.Mining.DomriaPrimaryMiner.Name(): domria.NewMiner(
			cfg.Mining.DomriaPrimaryMiner,
			db,
			gatherer,
			logger,
		),
		cfg.Mining.DomriaSecondaryMiner.Name(): domria.NewMiner(
			cfg.Mining.DomriaSecondaryMiner,
			db,
			gatherer,
			logger,
		),
	}
	job, ok := jobs[*alias]
	if !ok {
		_ = db.Close()
		logger.Fatal("main: mining failed to find the miner")
		return
	}
	miners := map[string]config.Miner{
		cfg.Mining.DomriaPrimaryMiner.Name():   cfg.Mining.DomriaPrimaryMiner,
		cfg.Mining.DomriaSecondaryMiner.Name(): cfg.Mining.DomriaSecondaryMiner,
	}
	miner := miners[*alias]
	if *isOnce {
		job.Run()
	} else {
		cron := gocron.New()
		if _, err = cron.AddJob(miner.Schedule(), job); err != nil {
			_ = db.Close()
			logger.Fatalf("main: mining failed to run, %v", err)
		}
		metrics.RunServer(miner.Metrics(), logger)
		cron.Run()
	}
	if err = database.CloseDatabase(db); err != nil {
		logger.Fatal(err)
	}
}
