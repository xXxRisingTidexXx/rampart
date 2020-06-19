package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining/config"
	"rampart/internal/mining/misc"
)

func NewProspector(housing misc.Housing, config *config.Domria) *Prospector {
	return &Prospector{
		housing,
		newFetcher(config.Fetcher),
		newSanitizer(config.Sanitizer),
		newGeocoder(config.Geocoder),
		newValidator(config.Validator),
	}
}

type Prospector struct {
	housing   misc.Housing
	fetcher   *fetcher
	sanitizer *sanitizer
	geocoder  *geocoder
	validator *validator
}

func (prospector *Prospector) Prospect() error {
	log.Debug("domria: prospector started")
	flats, err := prospector.fetcher.fetchFlats(prospector.housing)
	if err != nil {
		log.Debug("domria: prospector terminated")
		return err
	}
	flats = prospector.sanitizer.sanitizeFlats(flats)
	flats = prospector.geocoder.geocodeFlats(flats)
	flats = prospector.validator.validateFlats(flats)
	log.Debug("domria: prospector finished")
	return nil
}
