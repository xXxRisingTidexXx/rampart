package gauging

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"net/http"
)

func RunServer(config *config.Server, scheduler *Scheduler, logger log.FieldLogger) {
	server := &http.Server{
		Addr:           config.Address,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
		Handler:        &handler{scheduler, logger},
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("gauging: server met an error, %v", err)
		}
	}()
}
