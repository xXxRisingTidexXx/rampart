package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Geocoder struct {
	Timeout         misc.Timing       `yaml:"timeout"`
	Headers         map[string]string `yaml:"headers"`
	StatelessCities *misc.Set         `yaml:"statelessCities"`
	SearchURL       string            `yaml:"searchURL"`
	SRID            int               `yaml:"srid"`
}

func (geocoder *Geocoder) String() string {
	return fmt.Sprintf(
		"{%s %v %v %s %d}",
		geocoder.Timeout,
		geocoder.Headers,
		geocoder.StatelessCities,
		geocoder.SearchURL,
		geocoder.SRID,
	)
}
