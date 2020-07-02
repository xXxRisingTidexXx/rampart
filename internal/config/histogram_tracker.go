package config

import (
	"fmt"
)

type HistogramTracker struct {
	Name    string    `yaml:"name"`
	Help    string    `yaml:"help"`
	Buckets []float64 `yaml:"buckets"`
	Labels  []string  `yaml:"labels"`
}

func (tracker *HistogramTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Buckets, tracker.Labels)
}
