package mining

import (
	"encoding/json"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"net/http"
	"strings"
	"time"
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

func (a *geocodingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	if flat.HasLocation() {
		metrics.MessisGeocodings.WithLabelValues("located").Inc()
		return flat, nil
	}
	if flat.City == "" || flat.Street == "" || flat.HouseNumber == "" {
		metrics.MessisGeocodings.WithLabelValues("unlocated").Inc()
		return flat, nil
	}
	now := time.Now()
	positions, err := a.getPositions(flat)
	metrics.MessisGeocodingDuration.Observe(time.Since(now).Seconds())
	if err != nil {
		metrics.MessisGeocodings.WithLabelValues("failure").Inc()
		return flat, err
	}
	if len(positions) == 0 {
		metrics.MessisGeocodings.WithLabelValues("nothing").Inc()
		return flat, nil
	}
	metrics.MessisGeocodings.WithLabelValues("success").Inc()
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
		Miner:       flat.Miner,
		ParsingTime: flat.ParsingTime,
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
