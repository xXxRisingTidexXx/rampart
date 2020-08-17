package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func RunServer() {
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(9004),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1048576,
		Handler:        promhttp.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("metrics: server met an error, %v", err)
		}
	}()
}
