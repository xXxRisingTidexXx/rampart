package config

import (
	"fmt"
	"rampart/internal/misc"
)

type GeocodingTracker struct {
	Name       string           `yaml:"name"`
	Help       string           `yaml:"help"`
	Labels     []string         `yaml:"labels"`
	Categories *misc.Categories `yaml:"categories"`
}

func (tracker *GeocodingTracker) String() string {
	return fmt.Sprintf("{%s %s %v %v}", tracker.Name, tracker.Help, tracker.Labels, tracker.Categories)
}
