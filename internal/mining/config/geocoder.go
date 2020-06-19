package config

import (
	"fmt"
	"rampart/internal/mining/util"
)

type Geocoder struct {
	Timeout         util.Timeout      `yaml:"timeout"`
	Headers         map[string]string `yaml:"headers"`
	StatelessCities *util.Set         `yaml:"statelessCities"`
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
