package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
)

func NewProspector(housing mining.Housing, config *configs.Domria) mining.Prospector {
	return &prospector{housing, newFetcher(config.Fetcher), newValidator(config.Validator)}
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
		prospector.logFinish()
		return nil
	}
	_ = prospector.validator.validateFlats(flats)
	prospector.logFinish()
	return nil
}

func (prospector *prospector) logFinish() {
	log.Debugf("domria: %s housing prospector finished", prospector.housing)
}
