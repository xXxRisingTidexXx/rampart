package config

import (
	"fmt"
)

type Runners struct {
	Domria *Domria `yaml:"domria"`
}

func (runners *Runners) String() string {
	return fmt.Sprintf("{%v}", runners.Domria)
}
