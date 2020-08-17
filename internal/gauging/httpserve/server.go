package httpserve

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func RunServer() {
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(9003),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1048576,
		Handler:        newHandler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("httpserve: server met an error, %v", err)
		}
	}()
}
