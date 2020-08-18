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
		make(chan *dto.Flat, 600),
		misc.Set{"Київ": struct{}{}},
	}
	go gauger.run()
	return gauger
}

type Gauger struct {
	client       *http.Client
	flatChannel  chan *dto.Flat
	subwayCities misc.Set
}

func (gauger *Gauger) run() {
	for range gauger.flatChannel {

	}
}

func (gauger *Gauger) GaugeFlats(flats []*dto.Flat) {
	for _, flat := range flats {
		gauger.flatChannel <- flat
	}
}
