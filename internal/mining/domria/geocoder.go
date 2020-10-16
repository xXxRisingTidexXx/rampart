package domria

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
	"strings"
	"time"
)

// TODO: add new search endpoint at https://locationiq.com/ .
// TODO: check unstructured query.
func NewGeocoder(
	config config.Geocoder,
	drain *metrics.Drain,
	logger log.FieldLogger,
) *Geocoder {
	return &Geocoder{
		&http.Client{Timeout: config.Timeout},
		config.Headers,
		config.StatelessCities,
		config.SearchFormat,
		drain,
		logger,
	}
}

type Geocoder struct {
	client          *http.Client
	headers         misc.Headers
	statelessCities misc.Set
	searchFormat    string
	drain           *metrics.Drain
	logger          log.FieldLogger
}

func (geocoder *Geocoder) GeocodeFlats(flats []Flat) []Flat {
	newFlats := make([]Flat, 0, len(flats))
	for _, flat := range flats {
		if newFlat, ok := geocoder.geocodeFlat(flat); ok {
			newFlats = append(newFlats, newFlat)
		}
	}
	return newFlats
}

func (geocoder *Geocoder) geocodeFlat(flat Flat) (Flat, bool) {
	if flat.IsInspected && flat.IsLocated() {
		geocoder.drain.DrainNumber(metrics.LocatedGeocodingNumber)
		return flat, true
	}
	if !flat.IsAddressable() {
		geocoder.drain.DrainNumber(metrics.UnlocatedGeocodingNumber)
		return Flat{}, false
	}
	start := time.Now()
	positions, err := geocoder.getPositions(flat)
	geocoder.drain.DrainDuration(metrics.GeocodingDuration, start)
	if err != nil {
		geocoder.logger.WithFields(
			log.Fields{"source": flat.Source, "url": flat.URL},
		).Error(err)
		geocoder.drain.DrainNumber(metrics.FailedGeocodingNumber)
		return Flat{}, false
	}
	if len(positions) == 0 {
		geocoder.drain.DrainNumber(metrics.InconclusiveGeocodingNumber)
		return Flat{}, false
	}
	geocoder.drain.DrainNumber(metrics.SuccessfulGeocodingNumber)
	return Flat{
		Source:      flat.Source,
		URL:         flat.URL,
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
		Point:       orb.Point{float64(positions[0].Lon), float64(positions[0].Lat)},
		State:       flat.State,
		City:        flat.City,
		District:    flat.District,
		Street:      flat.Street,
		HouseNumber: flat.HouseNumber,
	}, true
}

func (geocoder *Geocoder) getPositions(flat Flat) ([]position, error) {
	state := ""
	if !geocoder.statelessCities.Contains(flat.City) {
		state = strings.ReplaceAll(flat.State, " ", "+")
	}
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			geocoder.searchFormat,
			state,
			strings.ReplaceAll(flat.City, " ", "+"),
			strings.ReplaceAll(flat.Street, " ", "+"),
			strings.ReplaceAll(flat.HouseNumber, " ", "+"),
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
	positions := make([]position, 0)
	if err = json.NewDecoder(response.Body).Decode(&positions); err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("domria: fetcher failed to unmarshal the positions, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("domria: geocoder failed to close the response body, %v", err)
	}
	return positions, nil
}
