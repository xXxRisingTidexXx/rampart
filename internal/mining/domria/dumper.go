package domria

import (
	gobytes "bytes"
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/dto"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

func NewDumper(config *config.Dumper) *Dumper {
	return &Dumper{&http.Client{Timeout: config.Timeout}, config.GaugingURL, config.Headers}
}

type Dumper struct {
	client     *http.Client
	gaugingURL string
	headers    misc.Headers
}

func (dumper *Dumper) DumpFlats(flats []*Flat) error {
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
		return fmt.Errorf("domria: dumper failed to marshal locations, %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, dumper.gaugingURL, gobytes.NewBuffer(bytes))
	if err != nil {
		return fmt.Errorf("domria: dumper failed to construct a request, %v", err)
	}
	dumper.headers.Inject(request)
	response, err := dumper.client.Do(request)
	if err != nil {
		return fmt.Errorf("domria: dumper failed to perform a request, %v", err)
	}
	if err = response.Body.Close(); err != nil {
		return fmt.Errorf("domria: dumper failed to close the response body, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("domria: dumper got response with status %s", response.Status)
	}
	return nil
}
