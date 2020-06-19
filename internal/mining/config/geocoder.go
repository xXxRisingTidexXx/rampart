package config

import (
	"fmt"
	"rampart/internal/mining/misc"
)

type Geocoder struct {
	Timeout         misc.Timeout      `yaml:"timeout"`
	Headers         map[string]string `yaml:"headers"`
	StatelessCities *misc.Set         `yaml:"statelessCities"`
	SearchURL       string            `yaml:"searchURL"`
	MinLookup       float64           `yaml:"minLookup"`
}

func (geocoder *Geocoder) String() string {
	return fmt.Sprintf(
		"{%s %v %v %s %.2f}",
		geocoder.Timeout,
		geocoder.Headers,
		geocoder.StatelessCities,
		geocoder.SearchURL,
		geocoder.MinLookup,
	)
}
