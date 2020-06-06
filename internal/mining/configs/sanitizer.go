package configs

import (
	"fmt"
)

type Sanitizer struct {
	OriginURLPrefix string `yaml:"originURLPrefix"`
	ImageURLPrefix  string `yaml:"imageURLPrefix"`
	StateEnding     string `yaml:"stateEnding"`
	StateSuffix     string `yaml:"stateSuffix"`
	DistrictEnding  string `yaml:"districtEnding"`
	DistrictSuffix  string `yaml:"districtSuffix"`
}

func (sanitizer *Sanitizer) String() string {
	return fmt.Sprintf(
		"{%s %s %s %s %s %s}",
		sanitizer.OriginURLPrefix,
		sanitizer.ImageURLPrefix,
		sanitizer.StateEnding,
		sanitizer.StateSuffix,
		sanitizer.DistrictEnding,
		sanitizer.DistrictSuffix,
	)
}
