package config

import (
	"fmt"
)

type CounterTracker struct {
	Name   string   `yaml:"name"`
	Help   string   `yaml:"help"`
	Labels []string `yaml:"labels"`
}

func (tracker *CounterTracker) String() string {
	return fmt.Sprintf("{%s %s %v}", tracker.Name, tracker.Help, tracker.Labels)
}
