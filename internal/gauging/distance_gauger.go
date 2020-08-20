package gauging

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"math"
	"net/http"
	gourl "net/url"
)

type distanceGauger struct {
	client         *http.Client
	headers        misc.Headers
	interpreterURL string
	minArea        float64
	noDistance     float64
}

func (gauger *distanceGauger) queryCollection(query string, params ...interface{}) (*geojson.FeatureCollection, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(gauger.interpreterURL, gourl.QueryEscape(fmt.Sprintf(query, params...))),
		nil,
	)
	if err != nil {
		return nil, err
	}
	gauger.headers.Inject(request)
	response, err := gauger.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("gauging: gauger got response with status %s", response.Status)
	}
	o := osm.OSM{}
	if err := xml.NewDecoder(response.Body).Decode(&o); err != nil {
		_ = response.Body.Close()
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	collection, err := osmgeojson.Convert(&o, osmgeojson.NoMeta(true))
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (gauger *distanceGauger) gaugeDistance(flat *dto.Flat, collection *geojson.FeatureCollection) float64 {
	geometries := make([]orb.Geometry, 0)
	for _, feature := range collection.Features {
		if planar.Area(feature.Geometry) >= gauger.minArea {
			geometries = append(geometries, feature.Geometry)
		}
	}
	if len(geometries) == 0 {
		return gauger.noDistance
	}
	point := orb.Point(flat.Point)
	distance := planar.DistanceFrom(geometries[0], point)
	for _, geometry := range geometries {
		distance = math.Min(distance, planar.DistanceFrom(geometry, point))
	}
	return distance
}
