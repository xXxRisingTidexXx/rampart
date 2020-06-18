package util

import (
	"gopkg.in/yaml.v3"
	"time"
)

type Timeout time.Duration

func (timeout *Timeout) UnmarshalYAML(node *yaml.Node) error {
	s := ""
	if err := node.Decode(&s); err != nil {
		return err
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*timeout = Timeout(duration)
	return nil
}

func (timeout Timeout) String() string {
	return time.Duration(timeout).String()
}
