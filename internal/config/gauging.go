package config

import (
	"fmt"
)

type Gauging struct {
	DSNParams     map[string]string `yaml:"dnsParams"`
	HTTPServer    *Server           `yaml:"httpServer"`
	MetricsServer *Server           `yaml:"metricsServer"`
}

func (gauging *Gauging) String() string {
	return fmt.Sprintf("{%v %v %v}", gauging.DSNParams, gauging.HTTPServer, gauging.MetricsServer)
}
