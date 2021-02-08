package mining

import (
	"github.com/paulmach/orb"
	"net/http"
	"time"
)

// TODO: should we add states?
func NewGeocodingAmplifier() Amplifier {
	return &geocodingAmplifier{
		&http.Client{Timeout: time.Second * 10},
		"https://nominatim.openstreetmap.org/search?state=%s&city=%s&street=%s+%s&format=json&countrycodes=ua",
	}
}

type geocodingAmplifier struct {
	client       *http.Client
	searchFormat string
}

// TODO: metrics.
func (a *geocodingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	if flat.Point.Lon() != 0 || flat.Point.Lat() != 0 {
		return flat, nil
	}
	if flat.City == "" || flat.Street == "" || flat.HouseNumber == "" {
		return flat, nil
	}
	positions, err := a.getPositions(flat)
	if err != nil {
		return Flat{}, err
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
	return nil, nil
}
