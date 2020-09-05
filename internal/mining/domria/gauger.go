package domria

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewGauger() *Gauger {
	return &Gauger{
		&http.Client{Timeout: 40 * time.Second},
		misc.Headers{"User-Agent": "rampart-mining-bot/1.0.0"},
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
