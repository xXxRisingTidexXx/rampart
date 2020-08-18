package domria

import (
	"database/sql"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"time"
)

func NewMiner(
	config *config.DomriaMiner,
	db *sql.DB,
	gatherer *metrics.Gatherer,
	logger *logging.Logger,
) *Miner {
	return &Miner{
		string(config.Housing),
		NewFetcher(config.Fetcher, gatherer),
		NewSanitizer(config.Sanitizer, gatherer),
		NewGeocoder(config.Geocoder, gatherer, logger),
		NewValidator(config.Validator, gatherer),
		NewStorer(config.Storer, db, gatherer, logger),
		gatherer,
		logger,
	}
}

type Miner struct {
	housing   string
	fetcher   *Fetcher
	sanitizer *Sanitizer
	geocoder  *Geocoder
	validator *Validator
	storer    *Storer
	gatherer  *metrics.Gatherer
	logger    *logging.Logger
}

func (miner *Miner) Run() {
	start := time.Now()
	if flats, err := miner.fetcher.FetchFlats(miner.housing); err != nil {
		miner.logger.Error(err)
	} else {
		flats = miner.sanitizer.SanitizeFlats(flats)
		flats = miner.geocoder.GeocodeFlats(flats)
		flats = miner.validator.ValidateFlats(flats)
		locations := miner.storer.StoreFlats(flats)

	}
	miner.gatherer.GatherTotalDuration(start)
	if err := miner.gatherer.Flush(); err != nil {
		miner.logger.Error(err)
	}
}
