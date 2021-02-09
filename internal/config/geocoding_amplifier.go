package config

import (
	"time"
)

type GeocodingAmplifier struct {
	Timeout      time.Duration `yaml:"timeout"`
	SearchFormat string        `yaml:"search-format"`
}
