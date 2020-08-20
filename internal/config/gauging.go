package config

import (
	"fmt"
)

type Gauging struct {
	DSNParams     map[string]string `yaml:"dsnParams"`
	Scheduler     *Scheduler        `yaml:"scheduler"`
	HTTPServer    *Server           `yaml:"httpServer"`
	MetricsServer *Server           `yaml:"metricsServer"`
}

func (gauging *Gauging) String() string {
	return fmt.Sprintf(
		"{%v %v %v %v}",
		gauging.DSNParams,
		gauging.Scheduler,
		gauging.HTTPServer,
		gauging.MetricsServer,
	)
}
