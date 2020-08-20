package gauging

import (
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewSubwayStationDistanceGauger(client *http.Client) Gauger {
	return &subwayStationDistanceGauger{
		distanceGauger{
			client,
			misc.Headers{},
			"",
			0,
			-1,
		},
		1200,
	}
}

type subwayStationDistanceGauger struct {
	distanceGauger
	searchRadius float64
}

func (gauger *subwayStationDistanceGauger) GaugeFlat(flat *dto.Flat) (float64, error) {
	point := orb.Point(flat.Point)
	collection, err := gauger.queryCollection(
		"node[station=subway](around:%f,%f,%f);out;",
		gauger.searchRadius,
		point.Lat(),
		point.Lon(),
	)
	if err != nil {
		return gauger.noDistance, err
	}
	return gauger.gaugeDistance(flat, collection), nil
}
