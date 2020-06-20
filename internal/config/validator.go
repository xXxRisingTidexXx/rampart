package config

import (
	"fmt"
)

type Validator struct {
	MinPrice        float64 `yaml:"minPrice"`
	MinTotalArea    float64 `yaml:"minTotalArea"`
	MaxTotalArea    float64 `yaml:"maxTotalArea"`
	MinLivingArea   float64 `yaml:"minLivingArea"`
	MinKitchenArea  float64 `yaml:"minKitchenArea"`
	MinRoomNumber   int     `yaml:"minRoomNumber"`
	MaxRoomNumber   int     `yaml:"maxRoomNumber"`
	MinSpecificArea float64 `yaml:"minSpecificArea"`
	MaxSpecificArea float64 `yaml:"maxSpecificArea"`
	MinFloor        int     `yaml:"minFloor"`
	MinTotalFloor   int     `yaml:"minTotalFloor"`
	MaxTotalFloor   int     `yaml:"maxTotalFloor"`
	MinLongitude    float64 `yaml:"minLongitude"`
	MaxLongitude    float64 `yaml:"maxLongitude"`
	MinLatitude     float64 `yaml:"minLatitude"`
	MaxLatitude     float64 `yaml:"maxLatitude"`
}

func (validator *Validator) String() string {
	return fmt.Sprintf(
		"{%.1f %.1f %.1f %.1f %.1f %d %d %.1f %.1f %d %d %d %.1f %.1f %.1f %.1f}",
		validator.MinPrice,
		validator.MinTotalArea,
		validator.MaxTotalArea,
		validator.MinLivingArea,
		validator.MinKitchenArea,
		validator.MinRoomNumber,
		validator.MaxRoomNumber,
		validator.MinSpecificArea,
		validator.MaxSpecificArea,
		validator.MinFloor,
		validator.MinTotalFloor,
		validator.MaxTotalFloor,
		validator.MinLongitude,
		validator.MaxLongitude,
		validator.MinLatitude,
		validator.MaxLatitude,
	)
}
