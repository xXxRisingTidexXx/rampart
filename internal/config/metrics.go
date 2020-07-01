package config

import (
	"fmt"
)

type Metrics struct {
	Server   *Server   `yaml:"server"`
	Gatherer *Gatherer `yaml:"gatherer"`
}

func (metrics *Metrics) String() string {
	return fmt.Sprintf("{%v %v}", metrics.Server, metrics.Gatherer)
}
