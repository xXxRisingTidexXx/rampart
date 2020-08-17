package gauging

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"net/http"
)

func RunServer(config *config.Server) {
	server := &http.Server{
		Addr:           config.Address,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
		Handler:        newHandler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("httpserve: server met an error, %v", err)
		}
	}()
}
