package main

import (
	"flag"
	_ "github.com/lib/pq"
	gocron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/database"
	"github.com/xXxRisingTidexXx/rampart/internal/mining"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/secrets"
)

// TODO: https://github.com/robertkrimen/otto - JS interpreter in golang to run https://github.com/tyrasd/osmtogeojson .
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
	short, gatherer := cfg.Mining.Miners, metrics.NewGatherer(*alias, db)
	miners := map[string]mining.Miner{
		short.DomriaPrimary.Alias:   domria.NewMiner(short.DomriaPrimary, db, gatherer, logger),
		short.DomriaSecondary.Alias: domria.NewMiner(short.DomriaSecondary, db, gatherer, logger),
	}
	miner := miners[*alias]
	if miner == nil {
		_ = db.Close()
		logger.Fatal("main: mining failed to find the miner")
		return
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
