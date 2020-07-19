package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Sanitizer struct {
	OriginURLPrefix         string            `yaml:"originURLPrefix"`
	ImageURLPrefix          string            `yaml:"imageURLPrefix"`
	StateDictionary         map[string]string `yaml:"stateDictionary"`
	StateSuffix             string            `yaml:"stateSuffix"`
	CityDictionary          map[string]string `yaml:"cityDictionary"`
	DistrictDictionary      map[string]string `yaml:"districtDictionary"`
	DistrictCitySwaps       *misc.Set         `yaml:"districtCitySwaps"`
	DistrictEnding          string            `yaml:"districtEnding"`
	DistrictSuffix          string            `yaml:"districtSuffix"`
	StreetReplacements      []string          `yaml:"streetReplacements"`
	HouseNumberReplacements []string          `yaml:"houseNumberReplacements"`
}

func (sanitizer *Sanitizer) String() string {
	return fmt.Sprintf(
		"{%s %s %v %s %v %v %v %s %s %v %v}",
		sanitizer.OriginURLPrefix,
		sanitizer.ImageURLPrefix,
		sanitizer.StateDictionary,
		sanitizer.StateSuffix,
		sanitizer.CityDictionary,
		sanitizer.DistrictDictionary,
		sanitizer.DistrictCitySwaps,
		sanitizer.DistrictEnding,
		sanitizer.DistrictSuffix,
		sanitizer.StreetReplacements,
		sanitizer.HouseNumberReplacements,
	)
}
