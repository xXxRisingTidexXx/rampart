package mining

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"math"
	"net/http"
	"net/url"
	"time"
)

// TODO: relative city center distance feature (with city diameter).
func NewGaugingAmplifier(config config.GaugingAmplifier) Amplifier {
	return &gaugingAmplifier{
		&http.Client{Timeout: config.Timeout},
		config.Host,
		config.InterpreterFormat,
		config.UserAgent,
		config.SubwayCities,
		-1,
		config.SSFSearchRadius,
		config.SSFMinDistance,
		config.SSFModifier,
		config.IZFSearchRadius,
		config.IZFMinArea,
		config.IZFMinDistance,
		config.IZFModifier,
		config.GZFSearchRadius,
		config.GZFMinArea,
		config.GZFMinDistance,
		config.GZFModifier,
	}
}

type gaugingAmplifier struct {
	client            *http.Client
	host              string
	interpreterFormat string
	userAgent         string
	subwayCities      misc.Set
	unknownDistance   float64
	ssfSearchRadius   float64
	ssfMinDistance    float64
	ssfModifier       float64
	izfSearchRadius   float64
	izfMinArea        float64
	izfMinDistance    float64
	izfModifier       float64
	gzfSearchRadius   float64
	gzfMinArea        float64
	gzfMinDistance    float64
	gzfModifier       float64
}

func (a *gaugingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	if !flat.HasLocation() {
		return flat, nil
	}
	ssf, err := a.gaugeSSF(flat)
	if err != nil {
		return flat, err
	}
	izf, err := a.gaugeIZF(flat)
	if err != nil {
		return flat, err
	}
	gzf, err := a.gaugeGZF(flat)
	if err != nil {
		return flat, err
	}
	return Flat{
		URL:         flat.URL,
		ImageURLs:   flat.ImageURLs,
		Price:       flat.Price,
		TotalArea:   flat.TotalArea,
		LivingArea:  flat.LivingArea,
		KitchenArea: flat.KitchenArea,
		RoomNumber:  flat.RoomNumber,
		Floor:       flat.Floor,
		TotalFloor:  flat.TotalFloor,
		Housing:     flat.Housing,
		Point:       flat.Point,
		City:        flat.City,
		Street:      flat.Street,
		HouseNumber: flat.HouseNumber,
		SSF:         ssf,
		IZF:         izf,
		GZF:         gzf,
		Miner:       flat.Miner,
		ParsingTime: flat.ParsingTime,
	}, nil
}

func (a *gaugingAmplifier) gaugeSSF(flat Flat) (float64, error) {
	if !a.subwayCities.Contains(flat.City) {
		metrics.MessisGaugings.WithLabelValues(a.host, "ssf", "subwayless").Inc()
		return 0, nil
	}
	now := time.Now()
	collection, err := a.queryOverpass(
		fmt.Sprintf(
			"node[station=subway](around:%f,%f,%f);out;",
			a.ssfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
		),
	)
	metrics.MessisGaugingDuration.WithLabelValues(a.host, "ssf").Observe(time.Since(now).Seconds())
	if err != nil {
		metrics.MessisGaugings.WithLabelValues(a.host, "ssf", "failure").Inc()
		return 0, err
	}
	ssf := 0.0
	for _, feature := range collection.Features {
		distance := a.gaugeGeoDistance(feature.Geometry, flat.Point)
		if distance != a.unknownDistance {
			ssf += 1 / math.Max(distance, a.ssfMinDistance)
		}
	}
	if ssf == 0 {
		metrics.MessisGaugings.WithLabelValues(a.host, "ssf", "nothing").Inc()
		return 0, nil
	}
	metrics.MessisGaugings.WithLabelValues(a.host, "ssf", "success").Inc()
	return ssf * a.ssfModifier, nil
}

func (a *gaugingAmplifier) queryOverpass(query string) (*geojson.FeatureCollection, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(a.interpreterFormat, a.host, url.QueryEscape(query)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", a.userAgent)
	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to make a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("mining: amplifier got response with status %s", response.Status)
	}
	var o osm.OSM
	if err := xml.NewDecoder(response.Body).Decode(&o); err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("mining: amplifier failed to unmarshal xml, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to close a response body, %v", err)
	}
	collection, err := osmgeojson.Convert(&o, osmgeojson.NoMeta(true))
	if err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to convert from osm to geojson, %v", err)
	}
	return collection, nil
}

func (a *gaugingAmplifier) gaugeGeoDistance(geometry orb.Geometry, point orb.Point) float64 {
	distance := a.unknownDistance
	switch geometry := geometry.(type) {
	case nil:
		return distance
	case orb.Point:
		return geo.DistanceHaversine(geometry, point)
	case orb.MultiPoint:
		return a.gaugeGeoDistanceToPoints(geometry, point)
	case orb.LineString:
		return a.gaugeGeoDistanceToPoints(geometry, point)
	case orb.MultiLineString:
		distance := a.unknownDistance
		for _, lineString := range geometry {
			newDistance := a.gaugeGeoDistanceToPoints(lineString, point)
			if a.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.Ring:
		return a.gaugeGeoDistanceToPoints(geometry, point)
	case orb.Polygon:
		for _, ring := range geometry {
			newDistance := a.gaugeGeoDistanceToPoints(ring, point)
			if a.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.MultiPolygon:
		for _, polygon := range geometry {
			for _, ring := range polygon {
				newDistance := a.gaugeGeoDistanceToPoints(ring, point)
				if a.isLower(newDistance, distance) {
					distance = newDistance
				}
			}
		}
		return distance
	case orb.Collection:
		for _, newGeometry := range geometry {
			newDistance := a.gaugeGeoDistance(newGeometry, point)
			if a.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.Bound:
		return a.gaugeGeoDistanceToPoints(geometry.ToRing(), point)
	default:
		return distance
	}
}

func (a *gaugingAmplifier) isLower(d1, d2 float64) bool {
	return d1 < d2 || d2 == a.unknownDistance && d2 < d1
}

func (a *gaugingAmplifier) gaugeGeoDistanceToPoints(points []orb.Point, point orb.Point) float64 {
	distance := a.unknownDistance
	for _, newPoint := range points {
		newDistance := geo.DistanceHaversine(newPoint, point)
		if a.isLower(newDistance, distance) {
			distance = newDistance
		}
	}
	return distance
}

func (a *gaugingAmplifier) gaugeIZF(flat Flat) (float64, error) {
	now := time.Now()
	collection, err := a.queryOverpass(
		fmt.Sprintf(
			"(way[landuse=industrial](around:%f,%f,%f);>;relation[landuse=industrial](around:%f,%"+
				"f,%f);>;);out;",
			a.izfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
			a.izfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
		),
	)
	metrics.MessisGaugingDuration.WithLabelValues(a.host, "izf").Observe(time.Since(now).Seconds())
	if err != nil {
		metrics.MessisGaugings.WithLabelValues(a.host, "izf", "failure").Inc()
		return 0, err
	}
	izf := 0.0
	for _, feature := range collection.Features {
		if area := geo.Area(feature.Geometry); area >= a.izfMinArea {
			distance := a.gaugeGeoDistance(feature.Geometry, flat.Point)
			if distance != a.unknownDistance {
				izf += area / math.Max(distance, a.izfMinDistance)
			}
		}
	}
	if izf == 0 {
		metrics.MessisGaugings.WithLabelValues(a.host, "izf", "nothing").Inc()
		return 0, nil
	}
	metrics.MessisGaugings.WithLabelValues(a.host, "izf", "success").Inc()
	return izf * a.izfModifier, nil
}

func (a *gaugingAmplifier) gaugeGZF(flat Flat) (float64, error) {
	now := time.Now()
	collection, err := a.queryOverpass(
		fmt.Sprintf(
			"(way[leisure=park](around:%f,%f,%f);>;relation[leisure=park](around:%f,%f,%f);>;);ou"+
				"t;",
			a.gzfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
			a.gzfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
		),
	)
	metrics.MessisGaugingDuration.WithLabelValues(a.host, "gzf").Observe(time.Since(now).Seconds())
	if err != nil {
		metrics.MessisGaugings.WithLabelValues(a.host, "gzf", "failure").Inc()
		return 0, err
	}
	gzf := 0.0
	for _, feature := range collection.Features {
		if area := geo.Area(feature.Geometry); area >= a.gzfMinArea {
			distance := a.gaugeGeoDistance(feature.Geometry, flat.Point)
			if distance != a.unknownDistance {
				gzf += area / math.Max(distance, a.gzfMinDistance)
			}
		}
	}
	if gzf == 0 {
		metrics.MessisGaugings.WithLabelValues(a.host, "gzf", "nothing").Inc()
		return 0, nil
	}
	metrics.MessisGaugings.WithLabelValues(a.host, "gzf", "success").Inc()
	return gzf * a.gzfModifier, nil
}
