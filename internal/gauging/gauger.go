package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewGauger() *Gauger {
	gauger := &Gauger{
		&http.Client{Timeout: 35 * time.Second},
		make(chan *intent, 600),
		misc.Set{"Київ": struct{}{}},
	}
	go gauger.run()
	return gauger
}

type Gauger struct {
	client        *http.Client
	intentChannel chan *intent
	subwayCities  misc.Set
}

func (gauger *Gauger) run() {
	for range gauger.intentChannel {

	}
}

func (gauger *Gauger) GaugeFlats(flats []*dto.Flat) {
	for _, flat := range flats {
		if gauger.subwayCities.Contains(flat.City) {
			gauger.intentChannel <- &intent{subwayStationDistance, flat}
		}
		gauger.intentChannel <- &intent{industrialZoneDistance, flat}
		gauger.intentChannel <- &intent{greenZoneDistance, flat}
	}
}
