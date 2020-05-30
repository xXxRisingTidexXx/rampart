package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
	"time"
)

func NewProspector() mining.Prospector {
	return &prospector{newFetcher(10, 10*time.Second)}
}

type prospector struct {
	fetcher *fetcher
}

func (prospector *prospector) Prospect(state, city string, housing mining.Housing) {
	log.Info(prospector.fetcher.fetchFlats(state, city, housing))
}
