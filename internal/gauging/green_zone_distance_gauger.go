package gauging

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/gauging/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewGreenZoneDistanceGauger(config *config.Gauger, client *http.Client) Gauger {
	return &greenZoneDistanceGauger{
		distanceGauger{
			client,
			config.Headers,
			config.InterpreterURL,
			config.MinArea,
		},
		config.SearchRadius,
	}
}

type greenZoneDistanceGauger struct {
	distanceGauger
	searchRadius float64
}

func (gauger *greenZoneDistanceGauger) GaugeFlat(flat *dto.Flat) (float64, error) {
	point, start := orb.Point(flat.Point), time.Now()
	collection, err := gauger.queryCollection(
		`(
		  way[leisure=park](around:%f,%f,%f);
		  >;
		  relation[leisure=park](around:%f,%f,%f);
		  >;
		);
		out;`,
		gauger.searchRadius,
		point.Lat(),
		point.Lon(),
		gauger.searchRadius,
		point.Lat(),
		point.Lon(),
	)
	metrics.GreenZoneDistanceDuration.WithLabelValues("gauging").Observe(time.Since(start).Seconds())
	if err != nil {
		metrics.GreenZoneDistanceGauging.WithLabelValues("failed").Inc()
		return misc.DistanceUnknown, fmt.Errorf(
			"gauging: green zone distance gauger failed to gauge flat, %v",
			err,
		)
	}
	distance, status := gauger.gaugeDistance(flat, collection), "successful"
	if distance == misc.DistanceUnknown {
		status = "inconclusive"
	}
	metrics.GreenZoneDistanceGauging.WithLabelValues(status).Inc()
	return distance, nil
}
