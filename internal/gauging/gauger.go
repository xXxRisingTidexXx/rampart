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
		time.Second,
		1500,
		0,
		2000,
		0.000004,
		1200,
		0.00001,
	}
	go gauger.run()
	return gauger
}

type Gauger struct {
	client                     *http.Client
	intentChannel              chan *intent
	subwayCities               misc.Set
	period                     time.Duration
	subwayStationSearchRadius  float64
	subwayStationMinArea       float64
	industrialZoneSearchRadius float64
	industrialZoneMinArea      float64
	greenZoneSearchRadius      float64
	greenZoneMinArea           float64
}

func (gauger *Gauger) run() {
	ticker := time.NewTicker(gauger.period)
	for intent := range gauger.intentChannel {
		<-ticker.C
		go gauger.gaugeIntent(intent)
	}
}

func (gauger *Gauger) gaugeIntent(intent *intent) {
	switch intent.target {
	case subwayStationDistance:

		break
	case industrialZoneDistance:

		break
	case greenZoneDistance:

		break
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
