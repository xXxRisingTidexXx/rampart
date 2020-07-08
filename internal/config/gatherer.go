package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Gatherer struct {
	RunTracker        *CounterTracker    `yaml:"runTracker"`
	Statuses          *misc.Statuses     `yaml:"statuses"`
	GeocodingTracker  *CounterTracker    `yaml:"geocodingTracker"`
	Categories        *misc.Categories   `yaml:"categories"`
	ValidationTracker *CounterTracker    `yaml:"validationTracker"`
	Verdicts          *misc.Verdicts     `yaml:"verdicts"`
	StoringTracker    *CounterTracker    `yaml:"storingTracker"`
	Consequences      *misc.Consequences `yaml:"consequences"`
	DurationTracker   *HistogramTracker  `yaml:"durationTracker"`
	Processes         *misc.Processes    `yaml:"processes"`
}

func (gatherer *Gatherer) String() string {
	return fmt.Sprintf(
		"{%v %v %v %v %v %v %v %v %v %v}",
		gatherer.RunTracker,
		gatherer.Statuses,
		gatherer.GeocodingTracker,
		gatherer.Categories,
		gatherer.ValidationTracker,
		gatherer.Verdicts,
		gatherer.StoringTracker,
		gatherer.Consequences,
		gatherer.DurationTracker,
		gatherer.Processes,
	)
}
