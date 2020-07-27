package domria

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/mining/metrics"
)

func NewValidator(config *config.Validator, gatherer *metrics.Gatherer) *Validator {
	return &Validator{
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
		gatherer,
	}
}

type Validator struct {
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
	gatherer        *metrics.Gatherer
}

func (validator *Validator) ValidateFlats(flats []*Flat) []*Flat {
	newFlats := make([]*Flat, 0, len(flats))
	for _, flat := range flats {
		if validator.validateFlat(flat) {
			newFlats = append(newFlats, flat)
			validator.gatherer.GatherApprovedValidation()
		} else {
			validator.gatherer.GatherDeniedValidation()
		}
	}
	return newFlats
}

//nolint:gocognit
func (validator *Validator) validateFlat(flat *Flat) bool {
	if flat.RoomNumber == 0 {
		return false
	}
	specificArea := flat.TotalArea / float64(flat.RoomNumber)
	return flat.OriginURL != "" &&
		validator.minPrice <= flat.Price &&
		validator.minTotalArea <= flat.TotalArea &&
		flat.TotalArea <= validator.maxTotalArea &&
		validator.minLivingArea <= flat.LivingArea &&
		flat.LivingArea < flat.TotalArea &&
		validator.minKitchenArea <= flat.KitchenArea &&
		flat.KitchenArea < flat.TotalArea &&
		validator.minRoomNumber <= flat.RoomNumber &&
		flat.RoomNumber <= validator.maxRoomNumber &&
		validator.minSpecificArea <= specificArea &&
		specificArea <= validator.maxSpecificArea &&
		validator.minFloor <= flat.Floor &&
		flat.Floor <= flat.TotalFloor &&
		validator.minTotalFloor <= flat.TotalFloor &&
		flat.TotalFloor <= validator.maxTotalFloor &&
		validator.minLongitude <= flat.Point.X() &&
		flat.Point.X() <= validator.maxLongitude &&
		validator.minLatitude <= flat.Point.Y() &&
		flat.Point.Y() <= validator.maxLatitude
}
