package domria

import (
	log "github.com/sirupsen/logrus"
	"rampart/internal/mining"
	"rampart/internal/mining/configs"
)

func newValidator(config *configs.Validator) *validator {
	return &validator{
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
		config.MinValidity,
	}
}

type validator struct {
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
	minValidity     float64
}

func (validator *validator) validateFlats(flats []*flat) []*flat {
	expectedLength := len(flats)
	if expectedLength == 0 {
		log.Debug("domria: validator skipped flats")
		return flats
	}
	newFlats := make([]*flat, 0, expectedLength)
	for _, flat := range flats {
		if validator.validateFlat(flat) {
			newFlats = append(newFlats, flat)
		}
	}
	actualLength := len(newFlats)
	log.Debugf("domria: validator approved %d flats", actualLength)
	if validity := float64(actualLength) / float64(expectedLength); validity < validator.minValidity {
		log.Warningf("domria: validator met low validity (%.3f)", validity)
	}
	return newFlats
}

//nolint:gocognit,gocyclo
func (validator *validator) validateFlat(flat *flat) bool {
	if flat == nil || flat.roomNumber == 0 {
		return false
	}
	specificArea := flat.totalArea / float64(flat.roomNumber)
	return flat.originURL != "" &&
		flat.updateTime != nil &&
		validator.minPrice <= flat.price &&
		validator.minTotalArea <= flat.totalArea &&
		flat.totalArea <= validator.maxTotalArea &&
		validator.minLivingArea <= flat.livingArea &&
		flat.livingArea < flat.totalArea &&
		validator.minKitchenArea <= flat.kitchenArea &&
		flat.kitchenArea < flat.totalArea &&
		validator.minRoomNumber <= flat.roomNumber &&
		flat.roomNumber <= validator.maxRoomNumber &&
		validator.minSpecificArea <= specificArea &&
		specificArea <= validator.maxSpecificArea &&
		validator.minFloor <= flat.floor &&
		flat.floor <= flat.totalFloor &&
		validator.minTotalFloor <= flat.totalFloor &&
		flat.totalFloor <= validator.maxTotalFloor &&
		(flat.housing == mining.Primary ||
			flat.housing == mining.Secondary) &&
		flat.state != "" &&
		flat.city != "" &&
		(flat.point == nil &&
			flat.district != "" &&
			flat.street != "" &&
			flat.houseNumber != "" ||
			flat.point != nil &&
				validator.minLongitude <= flat.point.X() &&
				flat.point.X() <= validator.maxLongitude &&
				validator.minLatitude <= flat.point.Y() &&
				flat.point.Y() <= validator.maxLatitude)
}
