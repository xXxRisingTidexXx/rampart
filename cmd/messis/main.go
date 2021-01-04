package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
)

func main() {
	isDebug := flag.Bool("debug", false, "Execute a single workflow instead of the whole schedule")
	alias := flag.String("miner", "", "Desired miner alias")
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithFields(log.Fields{"app": "messis", "miner": *alias})
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := sql.Open("postgres", c.Messis.DSN)
	if err != nil {
		entry.Fatalf("main: messis failed to open the db, %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		entry.Fatalf("main: messis failed to ping the db, %v", err)
	}
	drain := metrics.NewDrain(*alias, db, entry)
	jobs := map[string]cron.Job{
		c.Messis.DomriaPrimaryMiner.Name(): domria.NewMiner(
			c.Messis.DomriaPrimaryMiner,
			db,
			drain,
			entry,
		),
		c.Messis.DomriaSecondaryMiner.Name(): domria.NewMiner(
			c.Messis.DomriaSecondaryMiner,
			db,
			drain,
			entry,
		),
	}
	job, ok := jobs[*alias]
	if !ok {
		_ = db.Close()
		entry.Fatal("main: messis failed to find the miner")
		return
	}
	miners := map[string]config.Miner{
		c.Messis.DomriaPrimaryMiner.Name():   c.Messis.DomriaPrimaryMiner,
		c.Messis.DomriaSecondaryMiner.Name(): c.Messis.DomriaSecondaryMiner,
	}
	miner := miners[*alias]
	if *isDebug {
		job.Run()
	} else {
		scheduler := cron.New()
		if _, err = scheduler.AddJob(miner.Schedule(), job); err != nil {
			_ = db.Close()
			entry.Fatalf("main: messis failed to run, %v", err)
		}
		metrics.RunServer(miner.Metrics(), entry)
		scheduler.Run()
	}
	if err = db.Close(); err != nil {
		entry.Fatalf("main: messis failed to close the db, %v", err)
	}
}
