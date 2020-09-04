package main

import (
	"flag"
	gocron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/database"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func main() {
	isDebug := flag.Bool("debug", false, "Execute a single workflow instead of the whole schedule")
	alias := flag.String("miner", "", "Desired miner alias")
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithFields(log.Fields{"app": "mining", "miner": *alias})
	dsn, err := misc.GetEnv("RAMPART_DATABASE_DSN")
	if err != nil {
		entry.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := database.NewDatabase(dsn, cfg.Mining.DSNParams)
	if err != nil {
		entry.Fatal(err)
	}
	gatherer := metrics.NewGatherer(*alias, db)
	jobs := map[string]gocron.Job{
		cfg.Mining.DomriaPrimaryMiner.Name(): domria.NewMiner(
			cfg.Mining.DomriaPrimaryMiner,
			db,
			gatherer,
			entry,
		),
		cfg.Mining.DomriaSecondaryMiner.Name(): domria.NewMiner(
			cfg.Mining.DomriaSecondaryMiner,
			db,
			gatherer,
			entry,
		),
	}
	job, ok := jobs[*alias]
	if !ok {
		_ = db.Close()
		entry.Fatal("main: mining failed to find the miner")
		return
	}
	miners := map[string]config.Miner{
		cfg.Mining.DomriaPrimaryMiner.Name():   cfg.Mining.DomriaPrimaryMiner,
		cfg.Mining.DomriaSecondaryMiner.Name(): cfg.Mining.DomriaSecondaryMiner,
	}
	miner := miners[*alias]
	if *isDebug {
		job.Run()
	} else {
		cron := gocron.New()
		if _, err = cron.AddJob(miner.Schedule(), job); err != nil {
			_ = db.Close()
			entry.Fatalf("main: mining failed to run, %v", err)
		}
		metrics.RunServer(miner.Metrics(), entry)
		cron.Run()
	}
	if err = database.CloseDatabase(db); err != nil {
		entry.Fatal(err)
	}
}
