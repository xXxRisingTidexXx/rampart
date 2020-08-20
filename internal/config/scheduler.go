package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Scheduler struct {
	Timeout                      time.Duration `yaml:"timeout"`
	Capacity                     int           `yaml:"capacity"`
	Period                       time.Duration `yaml:"period"`
	SubwayCities                 misc.Set      `yaml:"subwayCities"`
	SubwayStationDistanceGauger  *Gauger       `yaml:"subwayStationDistanceGauger"`
	IndustrialZoneDistanceGauger *Gauger       `yaml:"industrialZoneDistanceGauger"`
	GreenZoneDistanceGauger      *Gauger       `yaml:"greenZoneDistanceGauger"`
}

func (scheduler *Scheduler) String() string {
	return fmt.Sprintf(
		"{%s %d %s %v %v %v %v}",
		scheduler.Timeout,
		scheduler.Capacity,
		scheduler.Period,
		scheduler.SubwayCities,
		scheduler.SubwayStationDistanceGauger,
		scheduler.IndustrialZoneDistanceGauger,
		scheduler.GreenZoneDistanceGauger,
	)
}
