package config

import (
	"fmt"
	"rampart/internal/misc"
)

type StoringTracker struct {
	Name         string             `yaml:"name"`
	Help         string             `yaml:"help"`
	Labels       []string           `yaml:"labels"`
	Consequences *misc.Consequences `yaml:"consequences"`
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
