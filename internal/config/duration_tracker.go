package config

import (
	"fmt"
)

type DurationTracker struct {
	Name      string     `yaml:"name"`
	Help      string     `yaml:"help"`
	Buckets   []float64  `yaml:"buckets"`
	Labels    []string   `yaml:"labels"`
	Processes *Processes `yaml:"processes"`
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
