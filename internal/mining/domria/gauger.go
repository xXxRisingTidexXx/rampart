package domria

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

func NewGauger(config *config.Gauger, gatherer *metrics.Gatherer, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: time.Duration(config.Timeout)},
		config.Headers,
		config.InterpreterURL,
		config.SearchRadius,
		-1,
		gatherer,
		logger,
	}
}

type Gauger struct {
	client         *http.Client
	headers        map[string]string
	interpreterURL string
	searchRadius   float64
	noDistance     float64
	gatherer       *metrics.Gatherer
	logger         log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, len(flats))
	for i, flat := range flats {
		_ = gauger.gaugeIndustrialZoneDistance(flat)
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

// TODO: add subwayless city optimization.
func (gauger *Gauger) gaugeSubwayStationDistance(flat *Flat) float64 {
	start := time.Now()
	gosm, err := gauger.queryOSM(
		"node[station=subway](around:%f,%f,%f);out;",
		gauger.searchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	gauger.gatherer.GatherSubwayGaugingDuration(start)
	if err != nil {
		gauger.gatherer.GatherFailedSubwayGauging()
		gauger.logger.WithFields(
			log.Fields{misc.FieldOriginURL: flat.OriginURL, misc.FieldSource: flat.Source},
		).Error(err)
		return gauger.noDistance
	}
	if len(gosm.Nodes) == 0 {
		gauger.gatherer.GatherInconclusiveSubwayGauging()
		return gauger.noDistance
	}
	distance := planar.Distance(flat.Point, gosm.Nodes[0].Point())
	for _, node := range gosm.Nodes {
		distance = math.Min(distance, planar.Distance(flat.Point, node.Point()))
	}
	gauger.gatherer.GatherSuccessfulSubwayGauging()
	return distance
}

func (gauger *Gauger) queryOSM(query string, params ...interface{}) (*osm.OSM, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(gauger.interpreterURL, url.QueryEscape(fmt.Sprintf(query, params...))),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to construct a request, %v", err)
	}
	for key, value := range gauger.headers {
		request.Header.Set(key, value)
	}
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
	return &gosm, nil
}

func (gauger *Gauger) log(flat *Flat, err error) {
	gauger.logger.WithFields(
		log.Fields{misc.FieldOriginURL: flat.OriginURL, misc.FieldSource: flat.Source},
	).Error(err)
}

func (gauger *Gauger) gaugeIndustrialZoneDistance(flat *Flat) float64 {
	collection, err := gauger.queryFeatureCollection(
		"(way[landuse=industrial](around:%f,%f,%f);relation[landuse=industrial](around:%f,%f,%f););out;",
		2000.0,
		flat.Point.Lat(),
		flat.Point.Lon(),
		2000.0,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	if err != nil {
		gauger.logger.WithFields(
			log.Fields{misc.FieldOriginURL: flat.OriginURL, misc.FieldSource: flat.Source},
		).Error(err)
		return gauger.noDistance
	}

	return 0
}

func (gauger *Gauger) queryFeatureCollection(
	query string,
	params ...interface{},
) (*geojson.FeatureCollection, error) {
	gosm, err := gauger.queryOSM(query, params...)
	if err != nil {
		return nil, err
	}
	collection, err := osmgeojson.Convert(
		gosm,
		osmgeojson.NoID(true),
		osmgeojson.NoMeta(true),
		osmgeojson.NoRelationMembership(true),
	)
	if err != nil {
		return nil, fmt.Errorf("domria: gauger failed to convert to geojson, %v", err)
	}
	return collection, nil
}
