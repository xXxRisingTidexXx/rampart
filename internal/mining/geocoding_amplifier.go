package mining

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"net/http"
	"strings"
)

// TODO: should we add states?
// TODO: LocationIQ geocoder.
func NewGeocodingAmplifier(config config.GeocodingAmplifier) Amplifier {
	return &geocodingAmplifier{
		&http.Client{Timeout: config.Timeout},
		config.SearchFormat,
		config.UserAgent,
	}
}

type geocodingAmplifier struct {
	client       *http.Client
	searchFormat string
	userAgent    string
}

// TODO: metrics.
func (a *geocodingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	if flat.HasLocation() {
		return flat, nil
	}
	if flat.City == "" || flat.Street == "" || flat.HouseNumber == "" {
		return flat, nil
	}
	positions, err := a.getPositions(flat)
	if err != nil {
		return flat, err
	}
	if len(positions) == 0 {
		return flat, nil
	}
	return Flat{
		URL:         flat.URL,
		ImageURLs:   flat.ImageURLs,
		Price:       flat.Price,
		TotalArea:   flat.TotalArea,
		LivingArea:  flat.LivingArea,
		KitchenArea: flat.KitchenArea,
		RoomNumber:  flat.RoomNumber,
		Floor:       flat.Floor,
		TotalFloor:  flat.TotalFloor,
		Housing:     flat.Housing,
		Point:       orb.Point{float64(positions[0].Lon), float64(positions[0].Lat)},
		City:        flat.City,
		Street:      flat.Street,
		HouseNumber: flat.HouseNumber,
	}, nil
}

func (a *geocodingAmplifier) getPositions(flat Flat) ([]position, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			a.searchFormat,
			strings.ReplaceAll(flat.City, " ", "+"),
			strings.ReplaceAll(flat.Street, " ", "+"),
			strings.ReplaceAll(flat.HouseNumber, " ", "+"),
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to construct a request, %v", err)
	}
	request.Header.Set("User-Agent", a.userAgent)
	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to make a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("mining: amplifier got response with status %s", response.Status)
	}
	positions := make([]position, 0)
	if err = json.NewDecoder(response.Body).Decode(&positions); err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("mining: amplifier failed to unmarshal positions, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return nil, fmt.Errorf("mining: amplifier failed to close a response body, %v", err)
	}
	return positions, nil
}
