package domria

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"net/http"
	"time"
)

func NewGauger(gatherer *metrics.Gatherer, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: 5 * time.Second},
		map[string]string{"User-Agent": "rampart/1.0"},
		"https://overpass.kumi.systems/api/interpreter?data=%s",
		gatherer,
		logger,
	}
}

type Gauger struct {
	client         *http.Client
	headers        map[string]string
	interpreterURL string
	gatherer       *metrics.Gatherer
	logger         log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	return flats
}
