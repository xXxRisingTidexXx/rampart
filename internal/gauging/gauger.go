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
	"time"
)

func NewGauger() *Gauger {
	gauger := &Gauger{
		&http.Client{Timeout: 35 * time.Second},
		misc.Headers{},
		"https://overpass-api.de/api/interpreter?data=%s",
		make(chan *intent, 600),
		misc.Set{"Київ": struct{}{}},
		time.Second,
		-1,
		1500,
		0,
		2000,
		0.000004,
		1200,
		0.00001,
	}
	go gauger.run()
	return gauger
}

type Gauger struct {
	client                     *http.Client
	headers                    misc.Headers
	interpreterURL             string
	intentChannel              chan *intent
	subwayCities               misc.Set
	period                     time.Duration
	noDistance                 float64
	subwayStationSearchRadius  float64
	subwayStationMinArea       float64
	industrialZoneSearchRadius float64
	industrialZoneMinArea      float64
	greenZoneSearchRadius      float64
	greenZoneMinArea           float64
}

func (gauger *Gauger) run() {
	ticker := time.NewTicker(gauger.period)
	for intent := range gauger.intentChannel {
		<-ticker.C
		go gauger.gaugeIntent(intent)
	}
}

func (gauger *Gauger) gaugeIntent(intent *intent) {
	switch intent.target {
	case subwayStationDistance:

		break
	case industrialZoneDistance:

		break
	case greenZoneDistance:

		break
	}
}

func (gauger *Gauger) GaugeFlats(flats []*dto.Flat) {
	for _, flat := range flats {
		if gauger.subwayCities.Contains(flat.City) {
			gauger.intentChannel <- &intent{subwayStationDistance, flat}
		}
		gauger.intentChannel <- &intent{industrialZoneDistance, flat}
		gauger.intentChannel <- &intent{greenZoneDistance, flat}
	}
}

func (gauger *Gauger) gaugeSubwayStationDistance(flat *dto.Flat) float64 {
	return gauger.noDistance
}

func (gauger *Gauger) query(query string, params ...interface{}) (*geojson.FeatureCollection, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(gauger.interpreterURL, gourl.QueryEscape(fmt.Sprintf(query, params...))),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("gauging: gauger failed to construct a request, %v", err)
	}
	gauger.headers.Inject(request)
	response, err := gauger.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("gauging: gauger failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("gauging: gauger got response with status %s", response.Status)
	}
	o := osm.OSM{}
	if err := xml.NewDecoder(response.Body).Decode(&o); err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("gauging: gauger failed to unmarshal the osm, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("gauging: gauger failed to close the response body, %v", err)
	}
	collection, err := osmgeojson.Convert(
		&o,
		osmgeojson.NoID(true),
		osmgeojson.NoMeta(true),
		osmgeojson.NoRelationMembership(true),
	)
	if err != nil {
		return nil, fmt.Errorf("gauging: gauger failed to convert to geojson, %v", err)
	}
	return collection, nil
}

func (gauger *Gauger) gaugeDistance(
	flat *dto.Flat,
	collection *geojson.FeatureCollection,
	minArea float64,
) float64 {
	geometries := make([]orb.Geometry, 0)
	for _, feature := range collection.Features {
		if planar.Area(feature.Geometry) >= minArea {
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