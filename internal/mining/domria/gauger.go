package domria

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb/planar"
	gosm "github.com/paulmach/osm"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

func NewGauger(gatherer *metrics.Gatherer, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: 5 * time.Second},
		map[string]string{"User-Agent": "rampart/1.0"},
		"https://overpass.kumi.systems/api/interpreter?data=%s",
		2000,
		gatherer,
		logger,
	}
}

type Gauger struct {
	client         *http.Client
	headers        map[string]string
	interpreterURL string
	searchRadius   float64
	gatherer       *metrics.Gatherer
	logger         log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, len(flats))
	for i, flat := range flats {
		if x := gauger.gaugeSubwayStationDistance(flat); x != -1 {
			gauger.logger.Info(flat.OriginURL, x)
		}
		newFlats[i] = flat
	}
	return newFlats
}

func (gauger *Gauger) gaugeSubwayStationDistance(flat *Flat) float64 {
	osm, err := gauger.query(
		"node[station=subway](around:%f,%f,%f);out;",
		gauger.searchRadius,
		flat.Point.Lat(),
		flat.Point.Lon(),
	)
	if err != nil {
		gauger.logger.WithField("origin_url", flat.OriginURL).Error(err)
		return -1
	}
	if len(osm.Nodes) == 0 {
		return -1
	}
	distance := planar.Distance(flat.Point, osm.Nodes[0].Point())
	for _, node := range osm.Nodes {
		distance = math.Min(distance, planar.Distance(flat.Point, node.Point()))
	}
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
