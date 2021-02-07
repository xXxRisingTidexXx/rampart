package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/mining"
	"io"
)

// TODO: graceful shutdown.
func main() {
	name := flag.String(
		"miner",
		"",
		"Set a concrete miner name to run it once; leave the field blank to up the whole messis",
	)
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "messis")
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
	miners := make(map[string]mining.Miner)
	for _, miner := range []mining.Miner{mining.NewDomriaMiner(c.Messis.DomriaMiner)} {
		miners[miner.Name()] = miner
	}
	if *name == "" {
		metrics.RunServer(c.Messis.Server, entry)
		scheduler := cron.New(cron.WithSeconds())
		for _, miner := range miners {
			entry := entry.WithField("miner", miner.Name())
			if _, err := scheduler.AddJob(miner.Spec(), wrap(miner, entry)); err != nil {
				entry.Errorf("main: messis failed to start miner, %v", err)
			}
		}
		scheduler.Run()
	} else {
		entry := entry.WithField("miner", *name)
		miner, ok := miners[*name]
		if !ok {
			entry.Fatal("main: messis failed to find the miner")
		}
		if flat, err := miner.MineFlat(); err != nil {
			entry.Error(err)
		} else {
			entry.Info(flat)
		}
	}
	if err := db.Close(); err != nil {
		entry.Fatalf("main: messis failed to close the db, %v", err)
	}
}

func wrap(miner mining.Miner, logger log.FieldLogger) cron.Job {
	return cron.FuncJob(
		func() {
			switch flat, err := miner.MineFlat(); err {
			case nil:
				logger.Info(flat)
			case io.EOF:
			default:
				logger.Error(err)
			}
		},
	)
}
