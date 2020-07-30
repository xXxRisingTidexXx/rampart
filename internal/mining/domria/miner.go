package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
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
		NewFetcher(config.Fetcher, gatherer),
		NewSanitizer(config.Sanitizer, gatherer),
		NewGeocoder(config.Geocoder, gatherer, logger),
		NewGauger(config.Gauger, gatherer, logger),
		NewValidator(config.Validator, gatherer),
		NewStorer(config.Storer, db, gatherer, logger),
		gatherer,
		logger,
	}
}

type Miner struct {
	housing   misc.Housing
	spec      string
	port      int
	fetcher   *Fetcher
	sanitizer *Sanitizer
	geocoder  *Geocoder
	gauger    *Gauger
	validator *Validator
	storer    *Storer
	gatherer  *metrics.Gatherer
	logger    log.FieldLogger
}

func (miner *Miner) Run() {
	start := time.Now()
	if flats, err := miner.fetcher.FetchFlats(miner.housing); err != nil {
		miner.logger.Error(err)
	} else {
		flats = miner.sanitizer.SanitizeFlats(flats)
		flats = miner.geocoder.GeocodeFlats(flats)
		flats = miner.gauger.GaugeFlats(flats)
		flats = miner.validator.ValidateFlats(flats)
		miner.storer.StoreFlats(flats)
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
