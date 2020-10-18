package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	gocron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/domria"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

// TODO: create photo & panorama table(s).
// TODO: relative city center distance feature (with city diameter).
// TODO: distance to workplace feature.
// TODO: shorten house number column (but research the actual width before).
func main() {
	isDebug := flag.Bool("debug", false, "Execute a single workflow instead of the whole schedule")
	alias := flag.String("miner", "", "Desired miner alias")
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithFields(log.Fields{"app": "mining", "miner": *alias})
	cfg, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	dsn, err := misc.GetEnv("RAMPART_DATABASE_DSN")
	if err != nil {
		entry.Fatal(err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		entry.Fatalf("main: mining failed to open the db, %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		entry.Fatalf("main: mining failed to ping the db, %v", err)
	}
	drain := metrics.NewDrain(*alias, db, entry)
	jobs := map[string]gocron.Job{
		cfg.Mining.DomriaPrimaryMiner.Name(): domria.NewMiner(
			cfg.Mining.DomriaPrimaryMiner,
			db,
			drain,
			entry,
		),
		cfg.Mining.DomriaSecondaryMiner.Name(): domria.NewMiner(
			cfg.Mining.DomriaSecondaryMiner,
			db,
			drain,
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
	if err = db.Close(); err != nil {
		entry.Fatalf("main: mining failed to close the db, %v", err)
	}
}
