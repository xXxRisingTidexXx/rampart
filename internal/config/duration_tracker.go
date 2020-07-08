package config

import (
	"fmt"
	"rampart/internal/misc"
)

type DurationTracker struct {
	Name      string          `yaml:"name"`
	Help      string          `yaml:"help"`
	Buckets   []float64       `yaml:"buckets"`
	Labels    []string        `yaml:"labels"`
	Processes *misc.Processes `yaml:"processes"`
}

func (tracker *DurationTracker) String() string {
	return fmt.Sprintf(
		"{%s %s %v %v %v}",
		tracker.Name,
		tracker.Help,
		tracker.Buckets,
		tracker.Labels,
		tracker.Processes,
	)
}
