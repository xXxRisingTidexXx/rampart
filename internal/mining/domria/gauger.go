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
		2000,
		25,
		1000,
		3000,
		35000,
		30,
		1000000000,
		2500,
		50000,
		20,
		0.001,
		gatherer,
		logger,
	}
}

type Gauger struct {
	client          *http.Client
	headers         misc.Headers
	interpreterURL  string
	subwayCities    misc.Set
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
	gatherer        *metrics.Gatherer
	logger          log.FieldLogger
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, len(flats))
	for i, flat := range flats {
		newFlats[i] = &Flat{
			Source:      flat.Source,
			OriginURL:   flat.OriginURL,
			ImageURL:    flat.ImageURL,
			MediaCount:  flat.MediaCount,
			UpdateTime:  flat.UpdateTime,
			IsInspected: flat.IsInspected,
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

func (gauger *Gauger) gaugeSSF(flat *Flat) float64 {
	if !gauger.subwayCities.Contains(flat.City) {
		return 0
	}
	collection, err := gauger.query("")
	if err != nil {
		gauger.logger.WithFields(
			log.Fields{"source": flat.Source, "origin_url": flat.OriginURL, "feature": "ssf"},
		).Errorf("domria: gauger failed to gauge, %v", err)
		return 0
	}
	return 1
}

func (gauger *Gauger) query(query string, params ...interface{}) (*geojson.FeatureCollection, error) {
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

func (gauger *Gauger) gaugeIZF(flat *Flat) float64 {
	return 0
}

func (gauger *Gauger) gaugeGZF(flat *Flat) float64 {
	return 0
}
