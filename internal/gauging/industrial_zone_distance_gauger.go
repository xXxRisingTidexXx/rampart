package gauging

import (
	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
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
			-1,
		},
		2000,
		0.001,
	}
}

type industrialZoneDistanceGauger struct {
	distanceGauger
	searchRadius float64
	minArea      float64
}

func (gauger *industrialZoneDistanceGauger) GaugeFlat(flat *dto.Flat) float64 {
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
		log.Error(err)
		return gauger.noDistance
	}
	return gauger.gaugeDistance(flat, collection, gauger.minArea)
}
