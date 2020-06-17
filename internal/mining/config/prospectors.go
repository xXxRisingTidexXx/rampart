package config

import (
	"fmt"
)

type Prospectors struct {
	Domria *Domria `yaml:"domria"`
}

func (prospectors *Prospectors) String() string {
	return fmt.Sprintf("{%v}", prospectors.Domria)
}
