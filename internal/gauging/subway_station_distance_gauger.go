package gauging

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewSubwayStationDistanceGauger(config *config.Gauger, client *http.Client) Gauger {
	return &subwayStationDistanceGauger{
		distanceGauger{
			client,
			config.Headers,
			config.InterpreterURL,
			config.MinArea,
		},
		config.SearchRadius,
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
		return misc.DistanceUnknown, fmt.Errorf(
			"gauging: subway station distance gauger failed to gauge flat, %v",
			err,
		)
	}
	return gauger.gaugeDistance(flat, collection), nil
}
