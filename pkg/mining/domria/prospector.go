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

func (prospector *prospector) Prospect(housing mining.Housing) error {
	flats, err := prospector.fetcher.fetchFlats(housing)
	log.Info(flats)
	return err
}
