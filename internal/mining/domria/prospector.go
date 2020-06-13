package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
)

func NewProspector(housing mining.Housing, config *configs.Domria) mining.Prospector {
	return &prospector{
		housing,
		newFetcher(config.Fetcher),
		newSanitizer(config.Sanitizer),
		newGeocoder(),
		newValidator(config.Validator),
	}
}

type prospector struct {
	housing   mining.Housing
	fetcher   *fetcher
	sanitizer *sanitizer
	geocoder  *geocoder
	validator *validator
}

func (prospector *prospector) Prospect() error {
	log.Debug("domria: prospector started")
	flats, err := prospector.fetcher.fetchFlats(prospector.housing)
	if err != nil {
		return err
	}
	flats = prospector.sanitizer.sanitizeFlats(flats)
	flats = prospector.geocoder.geocodeFlats(flats)
	flats = prospector.validator.validateFlats(flats)
	log.Debug("domria: prospector finished")
	return nil
}
