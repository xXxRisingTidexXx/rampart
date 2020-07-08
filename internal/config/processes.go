package config

import (
	"fmt"
)

type Processes struct {
	Fetching  string `yaml:"fetching"`
	Geocoding string `yaml:"geocoding"`
	Reading   string `yaml:"reading"`
	Creation  string `yaml:"creation"`
	Update    string `yaml:"update"`
	Run       string `yaml:"run"`
}

func (processes *Processes) String() string {
	return fmt.Sprintf(
		"{%s %s %s %s %s %s}",
		processes.Fetching,
		processes.Geocoding,
		processes.Reading,
		processes.Creation,
		processes.Update,
		processes.Run,
	)
}
