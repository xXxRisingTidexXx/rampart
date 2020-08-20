package gauging

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"net/http"
)

func NewIndustrialZoneDistanceGauger(config *config.Gauger, client *http.Client) Gauger {
	return &industrialZoneDistanceGauger{
		distanceGauger{
			client,
			config.Headers,
			config.InterpreterURL,
			config.MinArea,
			-1,
		},
		config.SearchRadius,
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
