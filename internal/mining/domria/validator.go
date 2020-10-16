package domria

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
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

func (validator *Validator) ValidateFlats(flats []Flat) []Flat {
	newFlats := make([]Flat, 0, len(flats))
	for _, flat := range flats {
		if validator.validateFlat(flat) {
			newFlats = append(newFlats, flat)
		}
	}
	return newFlats
}

func (validator *Validator) validateFlat(flat Flat) bool {
	if flat.RoomNumber == 0 {
		validator.drain.DrainNumber(metrics.DeniedValidationNumber)
		return false
	}
	specificArea := flat.TotalArea / float64(flat.RoomNumber)
	if flat.URL == "" ||
		validator.minPrice > flat.Price ||
		validator.minTotalArea > flat.TotalArea ||
		flat.TotalArea > validator.maxTotalArea ||
		validator.minLivingArea > flat.LivingArea ||
		flat.LivingArea >= flat.TotalArea ||
		validator.minKitchenArea > flat.KitchenArea ||
		flat.KitchenArea >= flat.TotalArea ||
		validator.minRoomNumber > flat.RoomNumber ||
		flat.RoomNumber > validator.maxRoomNumber ||
		validator.minSpecificArea > specificArea ||
		specificArea > validator.maxSpecificArea ||
		validator.minFloor > flat.Floor ||
		flat.Floor > flat.TotalFloor ||
		validator.minTotalFloor > flat.TotalFloor ||
		flat.TotalFloor > validator.maxTotalFloor ||
		validator.minLongitude > flat.Point.Lon() ||
		flat.Point.Lon() > validator.maxLongitude ||
		validator.minLatitude > flat.Point.Lat() ||
		flat.Point.Lat() > validator.maxLatitude ||
		!flat.IsLocated() {
		validator.drain.DrainNumber(metrics.DeniedValidationNumber)
		return false
	}
	if flat.MediaCount() < validator.minImageCount {
		validator.drain.DrainNumber(metrics.UninformativeValidationNumber)
		return false
	}
	if flat.IsSold {
		validator.drain.DrainNumber(metrics.SoldValidationNumber)
		return false
	}
	validator.drain.DrainNumber(metrics.ApprovedValidationNumber)
	return true
}
