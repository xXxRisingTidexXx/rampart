package config

import (
	"fmt"
)

type ValidationTracker struct {
	Name    string   `yaml:"name"`
	Help    string   `yaml:"help"`
	Labels  []string `yaml:"labels"`
	Verdict *Verdict `yaml:"verdict"`
}

func (tracker *ValidationTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Labels, tracker.Verdict)
}
