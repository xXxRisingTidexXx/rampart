package domria

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewGauger(config *config.Gauger, gatherer *metrics.Gatherer, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: config.Timeout},
		config.Headers,
		config.InterpreterURL,
		config.SubwayCities,
		gatherer,
		logger,
	}
}

type Gauger struct {
	client         *http.Client
	headers        misc.Headers
	interpreterURL string
	subwayCities   misc.Set
	gatherer       *metrics.Gatherer
	logger         log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	return flats
}
