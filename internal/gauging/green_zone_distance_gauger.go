package gauging

import (
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewGreenZoneDistanceGauger(client *http.Client) Gauger {
	return &greenZoneDistanceGauger{
		distanceGauger{
			client,
			misc.Headers{},
			"",
			-1,
		},
		1500,
		0.0004,
	}
}

type greenZoneDistanceGauger struct {
	distanceGauger
	searchRadius float64
	minArea      float64
}

func (gauger *greenZoneDistanceGauger) GaugeFlat(flat *dto.Flat) (float64, error) {
	point := orb.Point(flat.Point)
	collection, err := gauger.queryCollection(
		`(
		  way[leisure=park](around:%f,%f,%f);
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
	if err != nil {
		return gauger.noDistance, err
	}
	return gauger.gaugeDistance(flat, collection, gauger.minArea), nil // TODO: move min area to distance gauger.
}
