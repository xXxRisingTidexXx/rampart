package domria

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

func NewMiner(
	config config.DomriaMiner,
	db *sql.DB,
	drain *metrics.Drain,
	logger log.FieldLogger,
) *Miner {
	return &Miner{
		config.Housing,
		NewFetcher(config.Fetcher, drain, logger),
		NewSanitizer(config.Sanitizer, drain, logger),
		NewGeocoder(config.Geocoder, drain, logger),
		NewGauger(config.Gauger, drain, logger),
		NewValidator(config.Validator, drain),
		NewStorer(config.Storer, db, drain, logger),
		drain,
	}
}

type Miner struct {
	housing   misc.Housing
	fetcher   *Fetcher
	sanitizer *Sanitizer
	geocoder  *Geocoder
	gauger    *Gauger
	validator *Validator
	storer    *Storer
	drain     *metrics.Drain
}

func (miner *Miner) Run() {
	start := time.Now()
	flats := miner.fetcher.FetchFlats(miner.housing)
	flats = miner.sanitizer.SanitizeFlats(flats)
	flats = miner.geocoder.GeocodeFlats(flats)
	flats = miner.gauger.GaugeFlats(flats)
	flats = miner.validator.ValidateFlats(flats)
	miner.storer.StoreFlats(flats)
	miner.drain.DrainDuration(metrics.TotalDuration, start)
	miner.drain.Flush()
}
