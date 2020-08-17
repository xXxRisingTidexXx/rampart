package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"net/http"
)

func RunServer(config *config.Server, logger *logging.Logger) {
	server := &http.Server{
		Addr:           config.Address,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
		Handler:        promhttp.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("metrics: server met an error, %v", err)
		}
	}()
}
