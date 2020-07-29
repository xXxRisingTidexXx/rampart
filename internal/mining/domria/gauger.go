package domria

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"io/ioutil"
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
		if _, err := gauger.gaugeSubwayStationDistance(flat.Point); err != nil {
			gauger.logger.WithField("origin_url", flat.OriginURL).Error(err)
		}
		newFlats[i] = flat
	}
	return newFlats
}

func (gauger *Gauger) gaugeSubwayStationDistance(point *geom.Point) (float64, error) {
	bytes, err := gauger.query(
		"node[station=subway](around:%f,%f,%f);out;",
		gauger.searchRadius,
		point.Y(),
		point.X(),
	)
	gauger.logger.Info(string(bytes))
	return 0, err
}

func (gauger *Gauger) query(query string, parameters ...interface{}) ([]byte, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(gauger.interpreterURL, url.QueryEscape(fmt.Sprintf(query, parameters...))),
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
	if err = response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: gauger failed to close the response body, %v", err)
	}
	return bytes, nil
}
