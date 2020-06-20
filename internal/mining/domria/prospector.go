package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/misc"
)

func NewProspector(housing misc.Housing, config *config.Domria, db *sql.DB) *Prospector {
	return &Prospector{
		housing,
		newFetcher(config.Fetcher),
		newSanitizer(config.Sanitizer),
		newGeocoder(config.Geocoder),
		newValidator(config.Validator),
		newSifter(db, config.Sifter),
	}
}

type Prospector struct {
	housing   misc.Housing
	fetcher   *fetcher
	sanitizer *sanitizer
	geocoder  *geocoder
	validator *validator
	sifter    *sifter
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
	flats, err = prospector.sifter.siftFlats(flats)
	if err != nil {
		log.Debug("domria: prospector terminated")
		return err
	}
	log.Debug("domria: prospector finished")
	return nil
}
