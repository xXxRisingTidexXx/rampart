package main

import (
	"flag"
	_ "github.com/lib/pq"
	gocron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"rampart/internal/config"
	"rampart/internal/database"
	"rampart/internal/mining"
	"rampart/internal/mining/metrics"
	"rampart/internal/secrets"
	"syscall"
)

func main() {
	isOnce := flag.Bool("once", false, "Execute a single workflow instead of the whole schedule")
	alias := flag.String("miner", "", "Desired miner alias")
	flag.Parse()
	log.SetLevel(log.InfoLevel)
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
	gatherer := metrics.NewGatherer(*alias, cfg.Mining.Gatherer)
	miner, err := mining.FindMiner(*alias, cfg.Mining.Miners, db, gatherer)
	if err != nil {
		_ = db.Close()
		_ = gatherer.Unregister()
		log.Fatal(err)
	}
	if *isOnce {
		miner.Run()
	} else {
		cron := gocron.New()
		if _, err = cron.AddJob(miner.Spec(), miner); err != nil {
			_ = db.Close()
			_ = gatherer.Unregister()
			log.Fatalf("main: mining failed to schedule, %v", err)
		}
		metrics.RunServer(miner.Port(), cfg.Mining.Server)
		cron.Start()
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		<-signalChannel
		cron.Stop()
	}
	if err = gatherer.Unregister(); err != nil {
		_ = db.Close()
		log.Fatal(err)
	}
	if err = database.CloseDatabase(db); err != nil {
		log.Fatal(err)
	}
}
