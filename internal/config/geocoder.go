package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Geocoder struct {
	Timeout         Timing       `yaml:"timeout"`
	Headers         misc.Headers `yaml:"headers"`
	StatelessCities misc.Set     `yaml:"statelessCities"`
	SearchURL       string       `yaml:"searchURL"`
}

func (geocoder *Geocoder) String() string {
	return fmt.Sprintf(
		"{%s %v %v %s}",
		geocoder.Timeout,
		geocoder.Headers,
		geocoder.StatelessCities,
		geocoder.SearchURL,
	)
}
