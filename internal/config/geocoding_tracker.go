package config

import (
	"fmt"
)

type GeocodingTracker struct {
	Name       string      `yaml:"name"`
	Help       string      `yaml:"help"`
	Labels     []string    `yaml:"labels"`
	Categories *Categories `yaml:"categories"`
}

func (tracker *GeocodingTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Labels, tracker.Categories)
}
