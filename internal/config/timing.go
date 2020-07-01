package config

import (
	"gopkg.in/yaml.v3"
	"time"
)

type Timing time.Duration

func (timing *Timing) UnmarshalYAML(node *yaml.Node) error {
	s := ""
	if err := node.Decode(&s); err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*timing = Timing(duration)
	return nil
}

func (timing Timing) String() string {
	return time.Duration(timing).String()
}
