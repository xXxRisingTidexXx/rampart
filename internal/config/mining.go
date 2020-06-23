package config

import (
	"fmt"
)

type Mining struct {
	DSNParams map[string]string `yaml:"dsnParams"`
	Miners    *Miners           `yaml:"miners"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%v %v}", mining.DSNParams, mining.Miners)
}
