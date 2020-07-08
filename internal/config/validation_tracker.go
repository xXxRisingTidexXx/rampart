package config

import (
	"fmt"
)

type ValidationTracker struct {
	Name     string    `yaml:"name"`
	Help     string    `yaml:"help"`
	Labels   []string  `yaml:"labels"`
	Verdicts *Verdicts `yaml:"verdicts"`
}

func (tracker *ValidationTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Labels, tracker.Verdicts)
}
