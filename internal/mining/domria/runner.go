package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/misc"
	"time"
)

func NewRunner(housing misc.Housing, config *config.Domria, db *sql.DB) *Runner {
	return &Runner{
		housing,
		newFetcher(config.Fetcher),
		newSanitizer(config.Sanitizer),
		newGeocoder(config.Geocoder),
		newValidator(config.Validator),
		newSifter(db, config.Sifter),
		newStorer(db),
	}
}

type Runner struct {
	housing   misc.Housing
	fetcher   *fetcher
	sanitizer *sanitizer
	geocoder  *geocoder
	validator *validator
	sifter    *sifter
	storer    *storer
}

func (runner *Runner) Run() {
	start := time.Now()
	flats, err := runner.fetcher.fetchFlats(runner.housing)
	if err != nil {
		log.Error(err)
		return
	}
	flats = runner.sanitizer.sanitizeFlats(flats)
	flats = runner.geocoder.geocodeFlats(flats)
	flats = runner.validator.validateFlats(flats)
	if flats, err = runner.sifter.siftFlats(flats); err != nil {
		log.Error(err)
		return
	}
	if err = runner.storer.storeFlats(flats); err != nil {
		log.Error(err)
	} else {
		log.Debugf("domria: runner run (%.3fs)", time.Since(start).Seconds())
	}
}
