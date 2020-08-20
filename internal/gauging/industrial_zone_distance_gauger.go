package gauging

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewIndustrialZoneDistanceGauger(client *http.Client) Gauger {
	return &industrialZoneDistanceGauger{
		distanceGauger{
			client,
			misc.Headers{},
			"",
			0.0004,
			-1,
		},
		2000,
	}
}

type industrialZoneDistanceGauger struct {
	distanceGauger
	searchRadius float64
}

func (gauger *industrialZoneDistanceGauger) GaugeFlat(flat *dto.Flat) (float64, error) {
	point := orb.Point(flat.Point)
	collection, err := gauger.queryCollection(
		`(
		  way[landuse=industrial](around:%f,%f,%f);
		  relation[landuse=industrial](around:%f,%f,%f);
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
		return gauger.noDistance, fmt.Errorf(
			"gauging: industrial zone distance gauger failed to gauge flat, %v",
			err,
		)
	}
	return gauger.gaugeDistance(flat, collection), nil
}
