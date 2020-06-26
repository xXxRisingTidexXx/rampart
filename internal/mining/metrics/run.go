package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"rampart/internal/config"
	"strconv"
	"time"
)

func Run(config *config.Server) {
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(config.Port),
		ReadTimeout:    time.Duration(config.ReadTimeout),
		WriteTimeout:   time.Duration(config.WriteTimeout),
		MaxHeaderBytes: config.MaxHeaderBytes,
		Handler:        promhttp.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("metrics: server met an error, %v", err)
		}
	}()
}
