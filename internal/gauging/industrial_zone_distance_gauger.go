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

func NewIndustrialZoneDistanceGauger(config *config.Gauger, client *http.Client) Gauger {
	return &industrialZoneDistanceGauger{
		distanceGauger{
			client,
			config.Headers,
			config.InterpreterURL,
			config.MinArea,
		},
		config.SearchRadius,
	}
}

type industrialZoneDistanceGauger struct {
	distanceGauger
	searchRadius float64
}

func (gauger *industrialZoneDistanceGauger) GaugeFlat(flat *dto.Flat) (float64, error) {
	point, start := orb.Point(flat.Point), time.Now()
	collection, err := gauger.queryCollection(
		`(
		  way[landuse=industrial](around:%f,%f,%f);
		  relation[landuse=industrial](around:%f,%f,%f);
		  >;
		);
		out geom;`,
		gauger.searchRadius,
		point.Lat(),
		point.Lon(),
		gauger.searchRadius,
		point.Lat(),
		point.Lon(),
	)
	metrics.IndustrialZoneDistanceDuration.WithLabelValues("gauging").Observe(time.Since(start).Seconds())
	if err != nil {
		metrics.IndustrialZoneDistanceGauging.WithLabelValues("failed")
		return misc.DistanceUnknown, fmt.Errorf(
			"gauging: industrial zone distance gauger failed to gauge flat, %v",
			err,
		)
	}
	distance, status := gauger.gaugeDistance(flat, collection), "successful"
	if distance == misc.DistanceUnknown {
		status = "inconclusive"
	}
	metrics.IndustrialZoneDistanceGauging.WithLabelValues(status).Inc()
	return distance, nil
}
