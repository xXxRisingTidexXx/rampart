package config

import (
	"fmt"
)

type Metrics struct {
	Server *Server `yaml:"server"`
}

func (metrics *Metrics) String() string {
	return fmt.Sprintf("{%v}", metrics.Server)
}
