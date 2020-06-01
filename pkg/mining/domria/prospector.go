package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/pkg/mining"
	"rampart/pkg/mining/configs"
)

func NewProspector(housing mining.Housing, userAgent string, config *configs.Domria) mining.Prospector {
	return &prospector{housing, newFetcher(userAgent, config.Fetcher), newValidator(config.Validator)}
}

type prospector struct {
	housing   mining.Housing
	fetcher   *fetcher
	validator *validator
}

func (prospector *prospector) Prospect() error {
	log.Debugf("domria: %s housing prospector started", prospector.housing)
	flats, err := prospector.fetcher.fetchFlats(prospector.housing)
	if err != nil {
		return err
	}
	if len(flats) == 0 {
		return nil
	}
	flats = prospector.validator.validateFlats(flats)
	if len(flats) == 0 {
		return nil
	}
	log.Debugf("domria: %s housing prospector finished", prospector.housing)
	return nil
}
