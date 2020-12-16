package domria

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	gourl "net/url"
	"time"
)

// TODO: relative city center distance feature (with city diameter).
// TODO: distance to workplace feature.
func NewGauger(config config.Gauger, drain *metrics.Drain, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: config.Timeout},
		config.OverpassHosts,
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
		drain,
		logger,
	}
}

type Gauger struct {
	client          *http.Client
	overpassHosts   []string
	subwayCities    misc.Set
	unknownDistance float64
	ssfSearchRadius float64
	ssfMinDistance  float64
	ssfModifier     float64
	izfSearchRadius float64
	izfMinArea      float64
	izfMinDistance  float64
	izfModifier     float64
	gzfSearchRadius float64
	gzfMinArea      float64
	gzfMinDistance  float64
	gzfModifier     float64
	drain           *metrics.Drain
	logger          log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []Flat) []Flat {
	newFlats := make([]Flat, len(flats))
	for i, flat := range flats {
		newFlats[i] = Flat{
			Source:      flat.Source,
			URL:         flat.URL,
			Photos:      flat.Photos,
			Panoramas:   flat.Panoramas,
			UpdateTime:  flat.UpdateTime,
			Price:       flat.Price,
			TotalArea:   flat.TotalArea,
			LivingArea:  flat.LivingArea,
			KitchenArea: flat.KitchenArea,
			RoomNumber:  flat.RoomNumber,
			Floor:       flat.Floor,
			TotalFloor:  flat.TotalFloor,
			Housing:     flat.Housing,
			Complex:     flat.Complex,
			Point:       flat.Point,
			State:       flat.State,
			City:        flat.City,
			District:    flat.District,
			Street:      flat.Street,
			HouseNumber: flat.HouseNumber,
			SSF:         gauger.gaugeSSF(flat),
			IZF:         gauger.gaugeIZF(flat),
			GZF:         gauger.gaugeGZF(flat),
		}
	}
	return newFlats
}

func (gauger *Gauger) gaugeSSF(flat Flat) float64 {
	if !gauger.subwayCities.Contains(flat.City) {
		gauger.drain.DrainNumber(metrics.SubwaylessSSFGaugingNumber)
		return 0
	}
	entry := gauger.logger.WithFields(log.Fields{"url": flat.URL, "feature": "ssf"})
	start := time.Now()
	collection, err := gauger.queryOverpass(
		fmt.Sprintf(
			"node[station=subway](around:%f,%f,%f);out;",
			gauger.ssfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
		),
		entry,
	)
	gauger.drain.DrainDuration(metrics.SSFGaugingDuration, start)
	if err != nil {
		gauger.drain.DrainNumber(metrics.FailedSSFGaugingNumber)
		entry.Error(err)
		return 0
	}
	ssf := 0.0
	for _, feature := range collection.Features {
		distance := gauger.gaugeGeoDistance(feature.Geometry, flat.Point)
		if distance != gauger.unknownDistance {
			ssf += 1 / math.Max(distance, gauger.ssfMinDistance)
		}
	}
	if ssf == 0 {
		gauger.drain.DrainNumber(metrics.InconclusiveSSFGaugingNumber)
		return 0
	}
	gauger.drain.DrainNumber(metrics.SuccessfulSSFGaugingNumber)
	return ssf * gauger.ssfModifier
}

func (gauger *Gauger) queryOverpass(
	query string,
	logger log.FieldLogger,
) (*geojson.FeatureCollection, error) {
	data := gourl.QueryEscape(query)
	bytes, err := make([]byte, 0), io.EOF
	for i := 0; i < len(gauger.overpassHosts) && err != nil; i++ {
		if bytes, err = gauger.tryQuery(gauger.overpassHosts[i], data); err != nil {
			logger.WithField("host", gauger.overpassHosts[i]).Error(err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("domria: gauger exhausted hosts")
	}
	o := osm.OSM{}
	if err := xml.Unmarshal(bytes, &o); err != nil {
		return nil, fmt.Errorf("domria: gauger failed to unmarshal the xml, %v", err)
	}
	collection, err := osmgeojson.Convert(&o, osmgeojson.NoMeta(true))
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to convert from osm to geojson, %v", err)
	}
	return collection, nil
}

func (gauger *Gauger) tryQuery(host, data string) ([]byte, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://%s/api/interpreter?data=%s", host, data),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", misc.UserAgent)
	response, err := gauger.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: gauger got response with status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: gauger failed to read the response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: gauger failed to close the response body, %v", err)
	}
	return bytes, nil
}

func (gauger *Gauger) gaugeGeoDistance(geometry orb.Geometry, point orb.Point) float64 {
	distance := gauger.unknownDistance
	switch geometry := geometry.(type) {
	case nil:
		return distance
	case orb.Point:
		return geo.DistanceHaversine(geometry, point)
	case orb.MultiPoint:
		return gauger.gaugeGeoDistanceToPoints(geometry, point)
	case orb.LineString:
		return gauger.gaugeGeoDistanceToPoints(geometry, point)
	case orb.MultiLineString:
		distance := gauger.unknownDistance
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
		for _, ring := range geometry {
			newDistance := gauger.gaugeGeoDistanceToPoints(ring, point)
			if gauger.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.MultiPolygon:
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
		for _, newGeometry := range geometry {
			newDistance := gauger.gaugeGeoDistance(newGeometry, point)
			if gauger.isLower(newDistance, distance) {
				distance = newDistance
			}
		}
		return distance
	case orb.Bound:
		return gauger.gaugeGeoDistanceToPoints(geometry.ToRing(), point)
	default:
		return distance
	}
}

func (gauger *Gauger) isLower(d1, d2 float64) bool {
	return d1 < d2 || d2 == gauger.unknownDistance && d2 < d1
}

func (gauger *Gauger) gaugeGeoDistanceToPoints(points []orb.Point, point orb.Point) float64 {
	distance := gauger.unknownDistance
	for _, newPoint := range points {
		newDistance := geo.DistanceHaversine(newPoint, point)
		if gauger.isLower(newDistance, distance) {
			distance = newDistance
		}
	}
	return distance
}

func (gauger *Gauger) gaugeIZF(flat Flat) float64 {
	entry := gauger.logger.WithFields(log.Fields{"url": flat.URL, "feature": "izf"})
	start := time.Now()
	collection, err := gauger.queryOverpass(
		fmt.Sprintf(
			"(way[landuse=industrial](around:%f,%f,%f);>;relation[landuse=industrial](around:%f,%"+
				"f,%f);>;);out;",
			gauger.izfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
			gauger.izfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
		),
		entry,
	)
	gauger.drain.DrainDuration(metrics.IZFGaugingDuration, start)
	if err != nil {
		gauger.drain.DrainNumber(metrics.FailedIZFGaugingNumber)
		entry.Error(err)
		return 0
	}
	izf := 0.0
	for _, feature := range collection.Features {
		if area := geo.Area(feature.Geometry); area >= gauger.izfMinArea {
			distance := gauger.gaugeGeoDistance(feature.Geometry, flat.Point)
			if distance != gauger.unknownDistance {
				izf += area / math.Max(distance, gauger.izfMinDistance)
			}
		}
	}
	if izf == 0 {
		gauger.drain.DrainNumber(metrics.InconclusiveIZFGaugingNumber)
		return 0
	}
	gauger.drain.DrainNumber(metrics.SuccessfulIZFGaugingNumber)
	return izf * gauger.izfModifier
}

func (gauger *Gauger) gaugeGZF(flat Flat) float64 {
	entry := gauger.logger.WithFields(log.Fields{"url": flat.URL, "feature": "gzf"})
	start := time.Now()
	collection, err := gauger.queryOverpass(
		fmt.Sprintf(
			"(way[leisure=park](around:%f,%f,%f);>;relation[leisure=park](around:%f,%f,%f);>;);ou"+
				"t;",
			gauger.gzfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
			gauger.gzfSearchRadius,
			flat.Point.Lat(),
			flat.Point.Lon(),
		),
		entry,
	)
	gauger.drain.DrainDuration(metrics.GZFGaugingDuration, start)
	if err != nil {
		gauger.drain.DrainNumber(metrics.FailedGZFGaugingNumber)
		entry.Error(err)
		return 0
	}
	gzf := 0.0
	for _, feature := range collection.Features {
		if area := geo.Area(feature.Geometry); area >= gauger.gzfMinArea {
			distance := gauger.gaugeGeoDistance(feature.Geometry, flat.Point)
			if distance != gauger.unknownDistance {
				gzf += area / math.Max(distance, gauger.gzfMinDistance)
			}
		}
	}
	if gzf == 0 {
		gauger.drain.DrainNumber(metrics.InconclusiveGZFGaugingNumber)
		return 0
	}
	gauger.drain.DrainNumber(metrics.SuccessfulGZFGaugingNumber)
	return gzf * gauger.gzfModifier
}
