package gauging

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
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
		metrics.SubwayStationDistance.WithLabelValues("failed").Inc()
		return misc.DistanceUnknown, fmt.Errorf(
			"gauging: subway station distance gauger failed to gauge flat, %v",
			err,
		)
	}
	distance, category := gauger.gaugeDistance(flat, collection), "successful"
	if distance == misc.DistanceUnknown {
		category = "inconclusive"
	}
	metrics.SubwayStationDistance.WithLabelValues(category).Inc()
	return distance, nil
}
