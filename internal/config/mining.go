package config

import (
	"fmt"
)

type Mining struct {
	DSNParams   map[string]string `yaml:"dsnParams"`
	Prospectors *Prospectors      `yaml:"prospectors"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%v %v}", mining.DSNParams, mining.Prospectors)
}
