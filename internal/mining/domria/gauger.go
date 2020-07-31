package domria

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb/planar"
	gosm "github.com/paulmach/osm"
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
		}
	}
	return newFlats
}

func (gauger *Gauger) gaugeSubwayStationDistance(flat *Flat) float64 {
	start := time.Now()
	osm, err := gauger.query(
		"node[station=subway](around:%f,%f,%f);out;",
		gauger.searchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	gauger.gatherer.GatherSubwayGaugingDuration(start)
	if err != nil {
		gauger.gatherer.GatherFailedSubwayGauging()
		gauger.logger.WithField(misc.FieldOriginURL, flat.OriginURL).Error(err)
		return gauger.noDistance
	}
	if len(osm.Nodes) == 0 {
		gauger.gatherer.GatherInconclusiveSubwayGauging()
		return gauger.noDistance
	}
	distance := planar.Distance(flat.Point, osm.Nodes[0].Point())
	for _, node := range osm.Nodes {
		distance = math.Min(distance, planar.Distance(flat.Point, node.Point()))
	}
	gauger.gatherer.GatherSuccessfulSubwayGauging()
	return distance
}

func (gauger *Gauger) query(query string, params ...interface{}) (*gosm.OSM, error) {
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
	var osm gosm.OSM
	if err := xml.Unmarshal(bytes, &osm); err != nil {
		return nil, fmt.Errorf("domria: gauger failed to unmarshal the osm, %v", err)
	}
	return &osm, nil
}
