package logging

import (
	log "github.com/sirupsen/logrus"
)

func NewLogger(miner string) *Logger {
	return &Logger{log.WithField("miner", miner)}
}

type Logger struct {
	*log.Entry
}

func (logger *Logger) Problem(publication Publication, err error) {
	logger.WithFields(log.Fields{"url": publication.URL(), "body": publication.Body()}).Error(err)
}

func (logger *Logger) Inconsistency(message, url string, actual, expected interface{}) {
	logger.WithFields(log.Fields{"url": url, "actual": actual, "expected": expected}).Warning(message)
}
