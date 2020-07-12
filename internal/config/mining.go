package config

import (
	"fmt"
)

type Mining struct {
	DSNParams map[string]string `yaml:"dsnParams"`
	Server    *Server           `yaml:"server"`
	Miners    *Miners           `yaml:"miners"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%v %v %v}", mining.DSNParams, mining.Server, mining.Miners)
}
