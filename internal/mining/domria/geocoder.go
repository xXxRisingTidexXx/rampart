package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/logging"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func NewGeocoder(config *config.Geocoder, gatherer *metrics.Gatherer, logger *logging.Logger) *Geocoder {
	return &Geocoder{
		&http.Client{Timeout: config.Timeout},
		config.Headers,
		config.StatelessCities,
		config.SearchURL,
		gatherer,
		logger,
	}
}

type Geocoder struct {
	client          *http.Client
	headers         misc.Headers
	statelessCities misc.Set
	searchURL       string
	gatherer        *metrics.Gatherer
	logger          *logging.Logger
}

func (geocoder *Geocoder) GeocodeFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, 0, len(flats))
	for _, flat := range flats {
		if newFlat := geocoder.geocodeFlat(flat); newFlat != nil {
			newFlats = append(newFlats, newFlat)
		}
	}
	return newFlats
}

func (geocoder *Geocoder) geocodeFlat(flat *Flat) *Flat {
	if flat.Point.Lon() != 0 || flat.Point.Lat() != 0 {
		geocoder.gatherer.GatherLocatedGeocoding()
		return flat
	}
	if flat.State == "" ||
		flat.City == "" ||
		flat.District == "" ||
		flat.Street == "" ||
		flat.HouseNumber == "" {
		geocoder.gatherer.GatherUnlocatedGeocoding()
		return nil
	}
	start := time.Now()
	positions, err := geocoder.getLocations(flat)
	geocoder.gatherer.GatherGeocodingDuration(start)
	if err != nil {
		geocoder.logger.Problem(flat, err)
		geocoder.gatherer.GatherFailedGeocoding()
		return nil
	}
	if len(positions) == 0 {
		geocoder.gatherer.GatherInconclusiveGeocoding()
		return nil
	}
	geocoder.gatherer.GatherSuccessfulGeocoding()
	return &Flat{
		OriginURL:   flat.OriginURL,
		ImageURL:    flat.ImageURL,
		UpdateTime:  flat.UpdateTime,
		Price:       flat.Price,
		TotalArea:   flat.TotalArea,
		LivingArea:  flat.LivingArea,
		KitchenArea: flat.KitchenArea,
		RoomNumber:  flat.RoomNumber,
		Floor:       flat.Floor,
		TotalFloor:  flat.TotalFloor,
		Housing:     flat.Housing,
		Complex:     flat.Complex,
		Point:       orb.Point{float64(positions[0].Lon), float64(positions[0].Lat)},
		State:       flat.State,
		City:        flat.City,
		District:    flat.District,
		Street:      flat.Street,
		HouseNumber: flat.HouseNumber,
		Source:      flat.Source,
	}
}

func (geocoder *Geocoder) getLocations(flat *Flat) ([]*position, error) {
	whitespace, plus, state := " ", "+", ""
	if !geocoder.statelessCities.Contains(flat.City) {
		state = strings.ReplaceAll(flat.State, whitespace, plus)
	}
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			geocoder.searchURL,
			state,
			strings.ReplaceAll(flat.City, whitespace, plus),
			strings.ReplaceAll(flat.District, whitespace, plus),
			strings.ReplaceAll(flat.Street, whitespace, plus),
			strings.ReplaceAll(flat.HouseNumber, whitespace, plus),
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to construct a request, %v", err)
	}
	geocoder.headers.Inject(request)
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
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to close the response body, %v", err)
	}
	positions := make([]*position, 0)
	if err := json.Unmarshal(bytes, &positions); err != nil {
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the positions, %v", err)
	}
	return positions, nil
}
