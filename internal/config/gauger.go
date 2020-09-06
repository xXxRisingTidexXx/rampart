package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Gauger struct {
	Timeout         time.Duration `yaml:"timeout"`
	Headers         misc.Headers  `yaml:"headers"`
	InterpreterURL  string        `yaml:"interpreterURL"`
	SubwayCities    misc.Set      `yaml:"subwayCities"`
	SSFSearchRadius float64       `yaml:"ssfSearchRadius"`
	SSFMinDistance  float64       `yaml:"ssfMinDistance"`
	SSFModifier     float64       `yaml:"ssfModifier"`
	IZFSearchRadius float64       `yaml:"izfSearchRadius"`
	IZFMinArea      float64       `yaml:"izfMinArea"`
	IZFMinDistance  float64       `yaml:"izfMinDistance"`
	IZFModifier     float64       `yaml:"izfModifier"`
	GZFSearchRadius float64       `yaml:"gzfSearchRadius"`
	GZFMinArea      float64       `yaml:"gzfMinArea"`
	GZFMinDistance  float64       `yaml:"gzfMinDistance"`
	GZFModifier     float64       `yaml:"gzfModifier"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf(
		"{%s %v %s %v %.1f %.1f %.3f %.1f %.1f %.1f %.3f %.1f %.1f %.1f %.3f}",
		gauger.Timeout,
		gauger.Headers,
		gauger.InterpreterURL,
		gauger.SubwayCities,
		gauger.SSFSearchRadius,
		gauger.SSFMinDistance,
		gauger.SSFModifier,
		gauger.IZFSearchRadius,
		gauger.IZFMinArea,
		gauger.IZFMinDistance,
		gauger.IZFModifier,
		gauger.GZFSearchRadius,
		gauger.GZFMinArea,
		gauger.GZFMinDistance,
		gauger.GZFModifier,
	)
}
