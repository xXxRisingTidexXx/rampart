package gauging

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
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
	point := orb.Point(flat.Point)
	collection, err := gauger.queryCollection(
		`(
		  way[leisure=park](around:%f,%f,%f);
		  relation[leisure=park](around:%f,%f,%f);
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
	if err != nil {
		return misc.DistanceUnknown, fmt.Errorf(
			"gauging: green zone distance gauger failed to gauge flat, %v",
			err,
		)
	}
	return gauger.gaugeDistance(flat, collection), nil
}
