package config

import (
	"fmt"
)

type RunTracker struct {
	Name     string    `yaml:"name"`
	Help     string    `yaml:"help"`
	Labels   []string  `yaml:"labels"`
	Statuses *Statuses `yaml:"statuses"`
}

func (tracker *RunTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Labels, tracker.Statuses)
}
