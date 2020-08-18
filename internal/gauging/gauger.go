package gauging

import (
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"net/http"
	"time"
)

func NewGauger() *Gauger {
	return &Gauger{&http.Client{Timeout: 35 * time.Second}}
}

type Gauger struct {
	client *http.Client
}

func (gauger *Gauger) GaugeDistances(locations []*dto.Location) {

}
