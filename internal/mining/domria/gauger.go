package domria

import (
	gobytes "bytes"
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"time"
)

func NewGauger() *Gauger {
	return &Gauger{
		&http.Client{Timeout: 5 * time.Second},
		"http://rampart-gauging:9003",
		misc.Headers{"Content-Type": "application/json"},
	}
}

type Gauger struct {
	client     *http.Client
	gaugingURL string
	headers    misc.Headers
}

func (gauger *Gauger) GaugeFlats(flats []*Flat) error {
	length := len(flats)
	if length == 0 {
		return nil
	}
	locations := make([]*dto.Location, length)
	for i, flat := range flats {
		locations[i] = &dto.Location{OriginURL: flat.OriginURL, Point: geojson.Point(flat.Point)}
	}
	bytes, err := json.Marshal(locations)
	if err != nil {
		return fmt.Errorf("domria: gauger failed to marshal locations, %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, gauger.gaugingURL, gobytes.NewBuffer(bytes))
	if err != nil {
		return fmt.Errorf("domria: gauger failed to construct a request, %v", err)
	}
	gauger.headers.Inject(request)
	response, err := gauger.client.Do(request)
	if err != nil {
		return fmt.Errorf("domria: gauger failed to perform a request, %v", err)
	}
	if err = response.Body.Close(); err != nil {
		return fmt.Errorf("domria: gauger failed to close the response body, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("domria: gauger got response with status %s", response.Status)
	}
	return nil
}
