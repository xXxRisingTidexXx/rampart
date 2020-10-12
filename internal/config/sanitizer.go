package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Sanitizer struct {
	URLPrefix               string            `yaml:"url-prefix"`
	StateMap                map[string]string `yaml:"state-map"`
	StateSuffix             string            `yaml:"state-suffix"`
	CityMap                 map[string]string `yaml:"city-map"`
	DistrictMap             map[string]string `yaml:"district-map"`
	DistrictCitySwaps       misc.Set          `yaml:"district-city-swaps"`
	DistrictEnding          string            `yaml:"district-ending"`
	DistrictSuffix          string            `yaml:"district-suffix"`
	StreetReplacements      []string          `yaml:"street-replacements"`
	HouseNumberReplacements []string          `yaml:"house-number-replacements"`
	HouseNumberMaxLength    int               `yaml:"house-number-max-length"`
}
