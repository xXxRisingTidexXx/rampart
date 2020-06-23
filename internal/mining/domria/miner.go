package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/misc"
	"time"
)

func NewMiner(config *config.Domria, db *sql.DB) *Miner {
	return &Miner{
		config.Housing,
		config.Spec,
		newFetcher(config.Fetcher),
		newSanitizer(config.Sanitizer),
		newGeocoder(config.Geocoder),
		newValidator(config.Validator),
		newStorer(db, config.Storer),
	}
}

type Miner struct {
	housing   misc.Housing
	spec      string
	fetcher   *fetcher
	sanitizer *sanitizer
	geocoder  *geocoder
	validator *validator
	storer    *storer
}

func (miner *Miner) Run() {
	start := time.Now()
	flats, err := miner.fetcher.fetchFlats(miner.housing)
	if err != nil {
		log.Error(err)
		return
	}
	flats = miner.sanitizer.sanitizeFlats(flats)
	flats = miner.geocoder.geocodeFlats(flats)
	flats = miner.validator.validateFlats(flats)
	miner.storer.storeFlats(flats)
	log.Debugf("domria: miner run (%.3fs)", time.Since(start).Seconds())
}

func (miner *Miner) Spec() string {
	return miner.spec
}
