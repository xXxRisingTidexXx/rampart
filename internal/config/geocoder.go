package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Geocoder struct {
	Timeout         time.Duration `yaml:"timeout"`
	Headers         misc.Headers  `yaml:"headers"`
	StatelessCities misc.Set      `yaml:"stateless-cities"`
	SearchFormat    string        `yaml:"search-format"`
}
