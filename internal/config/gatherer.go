package config

import (
	"fmt"
)

type Gatherer struct {
	RunTracker        *RunTracker        `yaml:"runTracker"`
	GeocodingTracker  *GeocodingTracker  `yaml:"geocodingTracker"`
	ValidationTracker *ValidationTracker `yaml:"validationTracker"`
	StoringTracker    *StoringTracker    `yaml:"storingTracker"`
	DurationTracker   *DurationTracker   `yaml:"durationTracker"`
}

func (gatherer *Gatherer) String() string {
	return fmt.Sprintf(
		"{%v %v %v %v %v}",
		gatherer.RunTracker,
		gatherer.GeocodingTracker,
		gatherer.ValidationTracker,
		gatherer.StoringTracker,
		gatherer.DurationTracker,
	)
}
