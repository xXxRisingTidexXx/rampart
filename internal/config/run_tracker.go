package config

import (
	"fmt"
	"rampart/internal/misc"
)

type RunTracker struct {
	Name     string         `yaml:"name"`
	Help     string         `yaml:"help"`
	Labels   []string       `yaml:"labels"`
	Statuses *misc.Statuses `yaml:"statuses"`
}

func (tracker *RunTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Labels, tracker.Statuses)
}
