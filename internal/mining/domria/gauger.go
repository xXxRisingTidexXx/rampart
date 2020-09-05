package domria

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmgeojson"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	gourl "net/url"
)

func NewGauger(config *config.Gauger, gatherer *metrics.Gatherer, logger log.FieldLogger) *Gauger {
	return &Gauger{
		&http.Client{Timeout: config.Timeout},
		config.Headers,
		config.InterpreterURL,
		config.SubwayCities,
		gatherer,
		logger,
	}
}

type Gauger struct {
	client         *http.Client
	headers        misc.Headers
	interpreterURL string
	subwayCities   misc.Set
	gatherer       *metrics.Gatherer
	logger         log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	return flats
}

func (gauger *Gauger) queryCollection(query string, params ...interface{}) (*geojson.FeatureCollection, error) {
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
		return nil, fmt.Errorf("domria: gauger got response with status %s", response.Status)
	}
	o := osm.OSM{}
	if err := xml.NewDecoder(response.Body).Decode(&o); err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: gauger failed to unmarshal the xml, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: gauger failed to close the response body, %v", err)
	}
	collection, err := osmgeojson.Convert(&o, osmgeojson.NoMeta(true))
	if err != nil {
		return nil, err
	}
	return collection, nil
}
