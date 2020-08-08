package domria

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
	"math"
	"net/http"
	gourl "net/url"
	"time"
)

func NewGauger(config *config.Gauger, gatherer *metrics.Gatherer, logger *logging.Logger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: config.Timeout},
		config.Headers,
		config.InterpreterURL,
		config.NoDistance,
		config.SubwayCities,
		config.SubwayStationSearchRadius,
		config.IndustrialZoneSearchRadius,
		config.IndustrialZoneMinArea,
		config.GreenZoneSearchRadius,
		config.GreenZoneMinArea,
		gatherer,
		logger,
	}
}

type Gauger struct {
	client                     *http.Client
	headers                    misc.Headers
	interpreterURL             string
	noDistance                 float64
	subwayCities               misc.Set
	subwayStationSearchRadius  float64
	industrialZoneSearchRadius float64
	industrialZoneMinArea      float64
	greenZoneSearchRadius      float64
	greenZoneMinArea           float64
	gatherer                   *metrics.Gatherer
	logger                     *logging.Logger
}

// TODO: github.com/paulmach/osm can't parse some osm XMLs. Add
// TODO: metric to track such cases + think about possible
// TODO: workarounds. For instance, npm package.
func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, len(flats))
	for i, flat := range flats {
		newFlats[i] = &Flat{
			flat.OriginURL,
			flat.ImageURL,
			flat.UpdateTime,
			flat.Price,
			flat.TotalArea,
			flat.LivingArea,
			flat.KitchenArea,
			flat.RoomNumber,
			flat.Floor,
			flat.TotalFloor,
			flat.Housing,
			flat.Complex,
			flat.Point,
			gauger.gaugeSubwayStationDistance(flat),
			flat.State,
			flat.City,
			flat.District,
			flat.Street,
			flat.HouseNumber,
			flat.Source,
		}
	}
	return newFlats
}

func (gauger *Gauger) gaugeSubwayStationDistance(flat *Flat) float64 {
	if !gauger.subwayCities.Contains(flat.City) {
		gauger.gatherer.GatherAbsentSubwayGauging()
		return gauger.noDistance
	}
	start := time.Now()
	collection, err := gauger.query(
		"node[station=subway](around:%f,%f,%f);out;",
		gauger.subwayStationSearchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	gauger.gatherer.GatherSubwayGaugingDuration(start)
	if err != nil {
		gauger.gatherer.GatherFailedSubwayGauging()
		gauger.logger.Problem(flat, err)
		return gauger.noDistance
	}
	distance := gauger.gaugeDistance(flat, collection, 0)
	if distance == gauger.noDistance {
		gauger.gatherer.GatherInconclusiveSubwayGauging()
	} else {
		gauger.gatherer.GatherSuccessfulSubwayGauging()
	}
	return distance
}

func (gauger *Gauger) query(query string, params ...interface{}) (*geojson.FeatureCollection, error) {
	url := fmt.Sprintf(gauger.interpreterURL, gourl.QueryEscape(fmt.Sprintf(query, params...)))
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to construct a request, %v", err)
	}
	gauger.headers.Inject(request)
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
	gosm := osm.OSM{}
	if err := xml.Unmarshal(bytes, &gosm); err != nil {
		return nil, fmt.Errorf("domria: gauger failed to unmarshal the osm, %v", err)
	}
	collection, err := osmgeojson.Convert(
		&gosm,
		osmgeojson.NoID(true),
		osmgeojson.NoMeta(true),
		osmgeojson.NoRelationMembership(true),
	)
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to convert to geojson, %v", err)
	}
	actual, expected := len(collection.Features), len(gosm.Nodes)+len(gosm.Ways)+len(gosm.Relations)
	if actual != expected {
		gauger.logger.Inconsistency("domria: gauger poorly converted osm", url, actual, expected)
	}
	return collection, nil
}

func (gauger *Gauger) gaugeDistance(
	flat *Flat,
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
	distance := planar.DistanceFrom(geometries[0], flat.Point)
	for _, geometry := range geometries {
		distance = math.Min(distance, planar.DistanceFrom(geometry, flat.Point))
	}
	return distance
}

// TODO: inject metrics.
func (gauger *Gauger) gaugeIndustrialZoneDistance(flat *Flat) float64 {
	collection, err := gauger.query(
		`(
		  way[landuse=industrial](around:%f,%f,%f);
		  relation[landuse=industrial](around:%f,%f,%f);
		);
		out geom;`,
		gauger.industrialZoneSearchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
		gauger.industrialZoneSearchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	if err != nil {
		gauger.logger.Problem(flat, err)
		return gauger.noDistance
	}
	return gauger.gaugeDistance(flat, collection, gauger.industrialZoneMinArea)
}

// TODO: inject metrics.
func (gauger *Gauger) gaugeGreenZoneDistance(flat *Flat) float64 {
	collection, err := gauger.query(
		`(
		  way[leisure=park](around:%f,%f,%f);
		  relation[leisure=park](around:%f,%f,%f);
		);
		out geom;`,
		gauger.greenZoneSearchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
		gauger.greenZoneSearchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	if err != nil {
		gauger.logger.Problem(flat, err)
		return gauger.noDistance
	}
	return gauger.gaugeDistance(flat, collection, gauger.greenZoneMinArea)
}
