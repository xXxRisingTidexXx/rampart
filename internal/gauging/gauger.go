package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"net/http"
	"time"
)

func NewGauger() *Gauger {
	gauger := &Gauger{&http.Client{Timeout: 35 * time.Second}, make(chan *dto.Location, 200)}
	go gauger.run()
	return gauger
}

type Gauger struct {
	client          *http.Client
	locationChannel chan *dto.Location
}

func (gauger *Gauger) run() {
	for range gauger.locationChannel {

	}
}

func (gauger *Gauger) GaugeAmenities(locations []*dto.Location) {
	for _, location := range locations {
		gauger.locationChannel <- location
	}
}
