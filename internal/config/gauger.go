package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Gauger struct {
	Timeout                    Timing       `yaml:"timeout"`
	Headers                    misc.Headers `yaml:"headers"`
	InterpreterURL             string       `yaml:"interpreterURL"`
	NoDistance                 float64      `yaml:"noDistance"`
	SubwayCities               misc.Set     `yaml:"subwayCities"`
	SubwayStationSearchRadius  float64      `yaml:"subwayStationSearchRadius"`
	IndustrialZoneSearchRadius float64      `yaml:"industrialZoneSearchRadius"`
	IndustrialZoneMinArea      float64      `yaml:"industrialZoneMinArea"`
	GreenZoneSearchRadius      float64      `yaml:"greenZoneSearchRadius"`
	GreenZoneMinArea           float64      `yaml:"greenZoneMinArea"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf(
		"{%s %v %s %f %v %.1f %.1f %.9f %.1f %.9f}",
		gauger.Timeout,
		gauger.Headers,
		gauger.InterpreterURL,
		gauger.NoDistance,
		gauger.SubwayCities,
		gauger.SubwayStationSearchRadius,
		gauger.IndustrialZoneSearchRadius,
		gauger.IndustrialZoneMinArea,
		gauger.GreenZoneSearchRadius,
		gauger.GreenZoneMinArea,
	)
}
