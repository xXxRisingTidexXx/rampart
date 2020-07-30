package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func NewGeocoder(config *config.Geocoder, gatherer *metrics.Gatherer, logger log.FieldLogger) *Geocoder {
	return &Geocoder{
		&http.Client{Timeout: time.Duration(config.Timeout)},
		config.Headers,
		config.StatelessCities,
		config.SearchURL,
		" ",
		"+",
		gatherer,
		logger,
	}
}

type Geocoder struct {
	client          *http.Client
	headers         map[string]string
	statelessCities *misc.Set
	searchURL       string
	whitespace      string
	plus            string
	gatherer        *metrics.Gatherer
	logger          log.FieldLogger
}

func (geocoder *Geocoder) GeocodeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, 0, len(flats))
	for _, flat := range flats {
		if newFlat, err := geocoder.geocodeFlat(flat); err != nil {
			geocoder.logger.WithField("origin_url", flat.OriginURL).Error(err)
			geocoder.gatherer.GatherFailedGeocoding()
		} else if newFlat != nil {
			newFlats = append(newFlats, newFlat)
		}
	}
	return newFlats
}

func (geocoder *Geocoder) geocodeFlat(flat *Flat) (*Flat, error) {
	if flat.Point.Lon() != 0 || flat.Point.Lat() != 0 {
		geocoder.gatherer.GatherLocatedGeocoding()
		return flat, nil
	}
	if flat.State == "" ||
		flat.City == "" ||
		flat.District == "" ||
		flat.Street == "" ||
		flat.HouseNumber == "" {
		geocoder.gatherer.GatherUnlocatedGeocoding()
		return nil, nil
	}
	start := time.Now()
	bytes, err := geocoder.getLocations(flat)
	geocoder.gatherer.GatherGeocodingDuration(start)
	if err != nil {
		return nil, err
	}
	newFlat, err := geocoder.locateFlat(flat, bytes)
	if err != nil {
		return nil, err
	}
	if newFlat == nil {
		geocoder.gatherer.GatherInconclusiveGeocoding()
		return nil, nil
	}
	geocoder.gatherer.GatherSuccessfulGeocoding()
	return newFlat, nil
}

func (geocoder *Geocoder) getLocations(flat *Flat) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, geocoder.getSearchURL(flat), nil)
	if err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to construct a request, %v", err)
	}
	for key, value := range geocoder.headers {
		request.Header.Set(key, value)
	}
	response, err := geocoder.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to perform a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: geocoder got response with status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: geocoder failed to read the response body, %v", err)
	}
	if err = response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to close the response body, %v", err)
	}
	return bytes, nil
}

func (geocoder *Geocoder) getSearchURL(flat *Flat) string {
	state := ""
	if !geocoder.statelessCities.Contains(flat.City) {
		state = strings.ReplaceAll(flat.State, geocoder.whitespace, geocoder.plus)
	}
	return fmt.Sprintf(
		geocoder.searchURL,
		state,
		strings.ReplaceAll(flat.City, geocoder.whitespace, geocoder.plus),
		strings.ReplaceAll(flat.District, geocoder.whitespace, geocoder.plus),
		strings.ReplaceAll(flat.Street, geocoder.whitespace, geocoder.plus),
		strings.ReplaceAll(flat.HouseNumber, geocoder.whitespace, geocoder.plus),
	)
}

func (geocoder *Geocoder) locateFlat(flat *Flat, bytes []byte) (*Flat, error) {
	var locations []*location
	if err := json.Unmarshal(bytes, &locations); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the locations, %v", err)
	}
	if len(locations) == 0 {
		return nil, nil
	}
	return &Flat{
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
		orb.Point{float64(locations[0].Lon), float64(locations[0].Lat)},
		flat.State,
		flat.City,
		flat.District,
		flat.Street,
		flat.HouseNumber,
	}, nil
}
