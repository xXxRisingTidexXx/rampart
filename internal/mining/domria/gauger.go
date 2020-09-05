package domria

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewGauger(config *config.Gauger, gatherer *metrics.Gatherer, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: 40 * time.Second},
		config.Headers,
		"https://overpass.kumi.systems/api/interpreter?data=%s",
	}
}

type Gauger struct {
	client         *http.Client
	headers        misc.Headers
	interpreterURL string
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	return flats
}
