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
	"os"
	"os/signal"
	"syscall"
)

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
		// TODO: move to config.
		minings := make(chan mining.Flat, 100)
		scheduler := cron.New(cron.WithSeconds())
		for _, miner := range miners {
			entry := entry.WithField("miner", miner.Name())
			_, err := scheduler.AddJob(miner.Spec(), wrap(miner, minings, entry))
			if err != nil {
				entry.Fatalf("main: messis failed to start miner, %v", err)
			}
		}
		geocodings := make(chan mining.Flat, 100)
		go run(
			mining.NewGeocodingAmplifier(c.Messis.GeocodingAmplifier),
			minings,
			geocodings,
			entry.WithField("amplifier", "geocoding"),
		)
		gaugings := make(chan mining.Flat, 100)
		for _, amplifier := range c.Messis.GaugingAmplifiers {
			go run(
				mining.NewGaugingAmplifier(amplifier),
				geocodings,
				gaugings,
				entry.WithFields(log.Fields{"amplifier": "gauging", "host": amplifier.Host}),
			)
		}
		go run(
			mining.NewStoringAmplifier(c.Messis.StoringAmplifier, db),
			gaugings,
			nil,
			entry.WithField("amplifier", "storing"),
		)
		metrics.RunServer(c.Messis.Server, entry)
		scheduler.Start()
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		scheduler.Stop()
	} else {
		miner, ok := miners[*name]
		if !ok {
			entry.Fatalf("main: messis found no miner %s", *name)
		}
		if len(c.Messis.GaugingAmplifiers) == 0 {
			entry.Fatal("main: messis found no gauging amplifier")
		}
		flat, err := miner.MineFlat()
		if err != nil {
			entry.Fatal(err)
		}
		flat, err = mining.NewGeocodingAmplifier(c.Messis.GeocodingAmplifier).AmplifyFlat(flat)
		if err != nil {
			entry.Fatal(err)
		}
		flat, err = mining.NewGaugingAmplifier(c.Messis.GaugingAmplifiers[0]).AmplifyFlat(flat)
		if err != nil {
			entry.Fatal(err)
		}
		flat, err = mining.NewStoringAmplifier(c.Messis.StoringAmplifier, db).AmplifyFlat(flat)
		if err != nil {
			entry.Fatal(err)
		}
		entry.Info(flat)
	}
	if err := db.Close(); err != nil {
		entry.Fatalf("main: messis failed to close the db, %v", err)
	}
}

func wrap(miner mining.Miner, output chan<- mining.Flat, logger log.FieldLogger) cron.Job {
	return cron.FuncJob(
		func() {
			switch flat, err := miner.MineFlat(); err {
			case nil:
				output <- flat
			case io.EOF:
			default:
				logger.Error(err)
			}
		},
	)
}

func run(
	amplifier mining.Amplifier,
	input <-chan mining.Flat,
	output chan<- mining.Flat,
	logger log.FieldLogger,
) {
	for flat := range input {
		apartment, err := amplifier.AmplifyFlat(flat)
		if err != nil {
			logger.Error(err)
		} else if output != nil {
			output <- apartment
		}
	}
}
