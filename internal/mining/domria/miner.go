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
		NewGauger(),
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
	gauger    *Gauger
	gatherer  *metrics.Gatherer
	logger    *logging.Logger
}

func (miner *Miner) Run() {
	start := time.Now()
	flats, err := miner.fetcher.FetchFlats(miner.housing)
	if err != nil {
		miner.logger.Error(err)
	}
	flats = miner.sanitizer.SanitizeFlats(flats)
	flats = miner.geocoder.GeocodeFlats(flats)
	flats = miner.validator.ValidateFlats(flats)
	flats = miner.storer.StoreFlats(flats)
	if err := miner.gauger.GaugeFlats(flats); err != nil {
		miner.logger.Error(err)
	}
	miner.gatherer.GatherTotalDuration(start)
	if err := miner.gatherer.Flush(); err != nil {
		miner.logger.Error(err)
	}
}
