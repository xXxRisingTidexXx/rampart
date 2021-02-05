package domria

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
)

func NewValidator(config config.Validator, drain *metrics.Drain) *Validator {
	return &Validator{
		config.MinImageCount,
		config.MinPrice,
		config.MinTotalArea,
		config.MaxTotalArea,
		config.MinLivingArea,
		config.MinKitchenArea,
		config.MinRoomNumber,
		config.MaxRoomNumber,
		config.MinSpecificArea,
		config.MaxSpecificArea,
		config.MinFloor,
		config.MinTotalFloor,
		config.MaxTotalFloor,
		config.MinLongitude,
		config.MaxLongitude,
		config.MinLatitude,
		config.MaxLatitude,
		drain,
	}
}

type Validator struct {
	minImageCount   int
	minPrice        float64
	minTotalArea    float64
	maxTotalArea    float64
	minLivingArea   float64
	minKitchenArea  float64
	minRoomNumber   int
	maxRoomNumber   int
	minSpecificArea float64
	maxSpecificArea float64
	minFloor        int
	minTotalFloor   int
	maxTotalFloor   int
	minLongitude    float64
	maxLongitude    float64
	minLatitude     float64
	maxLatitude     float64
	drain           *metrics.Drain
}

// TODO: lock flats with empty city. Probably, we need reverse geocoding to avoid this shit.
func (v *Validator) ValidateFlats(flats []Flat) []Flat {
	newFlats := make([]Flat, 0, len(flats))
	for _, flat := range flats {
		if v.validateFlat(flat) {
			newFlats = append(newFlats, flat)
		}
	}
	return newFlats
}

func (v *Validator) validateFlat(flat Flat) bool {
	if flat.RoomNumber == 0 {
		v.drain.DrainNumber(metrics.DeniedValidationNumber)
		return false
	}
	specificArea := flat.TotalArea / float64(flat.RoomNumber)
	if flat.URL == "" ||
		v.minPrice > flat.Price ||
		v.minTotalArea > flat.TotalArea ||
		flat.TotalArea > v.maxTotalArea ||
		v.minLivingArea > flat.LivingArea ||
		flat.LivingArea >= flat.TotalArea ||
		v.minKitchenArea > flat.KitchenArea ||
		flat.KitchenArea >= flat.TotalArea ||
		v.minRoomNumber > flat.RoomNumber ||
		flat.RoomNumber > v.maxRoomNumber ||
		v.minSpecificArea > specificArea ||
		specificArea > v.maxSpecificArea ||
		v.minFloor > flat.Floor ||
		flat.Floor > flat.TotalFloor ||
		v.minTotalFloor > flat.TotalFloor ||
		flat.TotalFloor > v.maxTotalFloor ||
		v.minLongitude > flat.Point.Lon() ||
		flat.Point.Lon() > v.maxLongitude ||
		v.minLatitude > flat.Point.Lat() ||
		flat.Point.Lat() > v.maxLatitude ||
		!flat.IsLocated() {
		v.drain.DrainNumber(metrics.DeniedValidationNumber)
		return false
	}
	if flat.ImageCount() < v.minImageCount {
		v.drain.DrainNumber(metrics.UninformativeValidationNumber)
		return false
	}
	if flat.IsSold {
		v.drain.DrainNumber(metrics.SoldValidationNumber)
		return false
	}
	v.drain.DrainNumber(metrics.ApprovedValidationNumber)
	return true
}
