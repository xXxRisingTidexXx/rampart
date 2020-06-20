package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Sifter struct {
	TotalAreaBias  float64     `yaml:"totalAreaBias"`
	RoomNumberBias int         `yaml:"roomNumberBias"`
	FloorBias      int         `yaml:"floorBias"`
	TotalFloorBias int         `yaml:"totalFloorBias"`
	DistanceBias   float64     `yaml:"distanceBias"`
	UpdateTiming   misc.Timing `yaml:"updateTiming"`
}

func (sifter *Sifter) String() string {
	return fmt.Sprintf(
		"{%.1f %d %d %d %.9f %s}",
		sifter.TotalAreaBias,
		sifter.RoomNumberBias,
		sifter.FloorBias,
		sifter.TotalFloorBias,
		sifter.DistanceBias,
		sifter.UpdateTiming,
	)
}
