package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"net/http"
	"strconv"
	"time"
)

func RunServer(port int, config *config.Server, logger *logging.Logger) {
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		ReadTimeout:    time.Duration(config.ReadTimeout),
		WriteTimeout:   time.Duration(config.WriteTimeout),
		MaxHeaderBytes: config.MaxHeaderBytes,
		Handler:        promhttp.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("metrics: server met an error, %v", err)
		}
	}()
}
