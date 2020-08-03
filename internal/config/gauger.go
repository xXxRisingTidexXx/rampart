package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Gauger struct {
	Timeout                Timing            `yaml:"timeout"`
	Headers                map[string]string `yaml:"headers"`
	InterpreterURL         string            `yaml:"interpreterURL"`
	NoDistance             float64           `yaml:"noDistance"`
	SubwayCities           *misc.Set         `yaml:"subwayCities"`
	SubwaySearchRadius     float64           `yaml:"subwaySearchRadius"`
	IndustrialSearchRadius float64           `yaml:"industrialSearchRadius"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf(
		"{%s %v %s %f %v %.1f %.1f}",
		gauger.Timeout,
		gauger.Headers,
		gauger.InterpreterURL,
		gauger.NoDistance,
		gauger.SubwayCities,
		gauger.SubwaySearchRadius,
		gauger.IndustrialSearchRadius,
	)
}
