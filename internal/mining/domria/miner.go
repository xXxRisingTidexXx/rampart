package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"rampart/internal/config"
	"rampart/internal/mining/metrics"
	"rampart/internal/misc"
	"time"
)

func NewMiner(
	config *config.DomriaMiner,
	db *sql.DB,
	gatherer *metrics.Gatherer,
	logger log.FieldLogger,
) *Miner {
	return &Miner{
		config.Housing,
		config.Spec,
		config.Port,
		newFetcher(config.Fetcher, gatherer),
		newSanitizer(config.Sanitizer),
		newGeocoder(config.Geocoder, gatherer, logger),
		newValidator(config.Validator, gatherer),
		newStorer(db, gatherer, logger),
		gatherer,
		logger,
	}
}

type Miner struct {
	housing   misc.Housing
	spec      string
	port      int
	fetcher   *fetcher
	sanitizer *sanitizer
	geocoder  *geocoder
	validator *validator
	storer    *storer
	gatherer  *metrics.Gatherer
	logger    log.FieldLogger
}

func (miner *Miner) Run() {
	start := time.Now()
	if flats, err := miner.fetcher.fetchFlats(miner.housing); err != nil {
		miner.logger.Error(err)
	} else {
		flats = miner.sanitizer.sanitizeFlats(flats)
		flats = miner.geocoder.geocodeFlats(flats)
		flats = miner.validator.validateFlats(flats)
		miner.storer.storeFlats(flats)
		miner.gatherer.GatherSuccess()
	}
	miner.gatherer.GatherTotalDuration(start)
	if err := miner.gatherer.Flush(); err != nil {
		miner.logger.Error(err)
	}
}

func (miner *Miner) Spec() string {
	return miner.spec
}

func (miner *Miner) Port() int {
	return miner.port
}
