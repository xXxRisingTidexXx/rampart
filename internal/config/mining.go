package config

import (
	"fmt"
)

type Mining struct {
	DSNParams map[string]string `yaml:"dsnParams"`
	Metrics   *Metrics          `yaml:"metrics"`
	Miners    *Miners           `yaml:"miners"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%v %v %v}", mining.DSNParams, mining.Metrics, mining.Miners)
}
