package gauging

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	gourl "net/url"
)

type distanceGauger struct {
	client         *http.Client
	headers        misc.Headers
	interpreterURL string
	minArea        float64
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

func (gauger *distanceGauger) gaugeDistance(
	flat *dto.Flat,
	collection *geojson.FeatureCollection,
) float64 {
	features := make([]*geojson.Feature, 0)
	for _, feature := range collection.Features {
		if geo.Area(feature.Geometry) >= gauger.minArea {
			features = append(features, feature)
		}
	}
	point, id, distance := orb.Point(flat.Point), "", misc.DistanceUnknown
	for _, feature := range features {
		newDistance := gauger.gaugeGeoDistance(feature.Geometry, point)
		if gauger.isLower(newDistance, distance) {
			if newID, ok := feature.ID.(string); ok {
				id = newID
			} else {
				id = ""
			}
			distance = newDistance
		}
	}
	log.WithFields(log.Fields{"point": point, "id": id, "distance": distance}).Info("gauging: success")
	return distance
}

func (gauger *distanceGauger) isLower(distance1, distance2 float64) bool {
	return distance1 < distance2 || distance2 == misc.DistanceUnknown && distance2 < distance1
}

func (gauger *distanceGauger) gaugeGeoDistance(geometry orb.Geometry, point orb.Point) float64 {
	if geometry == nil {
		return misc.DistanceUnknown
	}
	switch geometry := geometry.(type) {
	case orb.Point:
		return geo.DistanceHaversine(geometry, point)
	case orb.MultiPoint:
		return gauger.gaugeGeoDistanceToPoints(geometry, point)
	case orb.LineString:
		return gauger.gaugeGeoDistanceToPoints(geometry, point)
	case orb.MultiLineString:
		distance := misc.DistanceUnknown
		for _, lineString := range geometry {
			newDistance := gauger.gaugeGeoDistanceToPoints(lineString, point)
			if gauger.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.Ring:
		return gauger.gaugeGeoDistanceToPoints(geometry, point)
	case orb.Polygon:
		distance := misc.DistanceUnknown
		for _, ring := range geometry {
			newDistance := gauger.gaugeGeoDistanceToPoints(ring, point)
			if gauger.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.MultiPolygon:
		distance := misc.DistanceUnknown
		for _, polygon := range geometry {
			for _, ring := range polygon {
				newDistance := gauger.gaugeGeoDistanceToPoints(ring, point)
				if gauger.isLower(newDistance, distance) {
					distance = newDistance
				}
			}
		}
		return distance
	case orb.Collection:
		distance := misc.DistanceUnknown
		for _, newGeometry := range geometry {
			newDistance := gauger.gaugeGeoDistance(newGeometry, point)
			if gauger.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.Bound:
		return gauger.gaugeGeoDistanceToPoints(geometry.ToRing(), point)
	}
	return misc.DistanceUnknown
}

func (gauger *distanceGauger) gaugeGeoDistanceToPoints(points []orb.Point, point orb.Point) float64 {
	distance := misc.DistanceUnknown
	for _, newPoint := range points {
		newDistance := geo.DistanceHaversine(newPoint, point)
		if gauger.isLower(newDistance, distance) {
			distance = newDistance
		}
	}
	return distance
}
