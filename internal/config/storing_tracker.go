package config

import (
	"fmt"
)

type StoringTracker struct {
	Name         string        `yaml:"name"`
	Help         string        `yaml:"help"`
	Labels       []string      `yaml:"labels"`
	Consequences *Consequences `yaml:"consequences"`
}

func (tracker *StoringTracker) String() string {
	return fmt.Sprintf(
		"{%s %s %v %v}",
		tracker.Name,
		tracker.Help,
		tracker.Labels,
		tracker.Consequences,
	)
}
