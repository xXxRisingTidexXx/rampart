package config

import (
	"fmt"
)

type Mining struct {
	DSNParams map[string]string `yaml:"dsnParams"`
	Runners   *Runners          `yaml:"runners"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%v %v}", mining.DSNParams, mining.Runners)
}
