package gauging

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewSubwayStationDistanceGauger(client *http.Client, ) Gauger {
	return &subwayStationDistanceGauger{
		distanceGauger{
			client,
			misc.Headers{},
			"",
			-1,
		},
		1200,
	}
}

type subwayStationDistanceGauger struct {
	distanceGauger
	searchRadius float64
}

func (gauger *subwayStationDistanceGauger) GaugeFlat(flat *dto.Flat) float64 {
	collection, err := gauger.queryCollection(``)
	if err != nil {
		log.Error(err)
		return gauger.noDistance
	}
	return gauger.gaugeDistance(flat, collection, 0)
}
