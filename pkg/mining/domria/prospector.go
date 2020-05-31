package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
)

func NewProspector(housing mining.Housing) mining.Prospector {
	return &prospector{housing, newFetcher()}
}

type prospector struct {
	housing mining.Housing
	fetcher *fetcher
}

func (prospector *prospector) Prospect() error {
	log.Debugf("domria: %s housing prospector started", prospector.housing)
	_, err := prospector.fetcher.fetchFlats(prospector.housing)
	if err != nil {
		return err
	}
	log.Debugf("domria: %s housing prospector finished", prospector.housing)
	return nil
}
