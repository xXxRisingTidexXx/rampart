package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/database"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/secrets"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	scr, err := secrets.NewSecrets()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewDatabase(scr.DSN, cfg.Gauging.DSNParams)
	if err != nil {
		log.Fatal(err)
	}
	gauging.RunServer(cfg.Gauging.HTTPServer, gauging.NewScheduler(db))
	metrics.RunServer(cfg.Gauging.MetricsServer)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
	if err = database.CloseDatabase(db); err != nil {
		log.Fatal(err)
	}
}
