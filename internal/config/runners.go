package config

import (
	"fmt"
)

type Runners struct {
	DomriaPrimary   *Domria `yaml:"domriaPrimary"`
	DomriaSecondary *Domria `yaml:"domriaSecondary"`
}

func (runners *Runners) String() string {
	return fmt.Sprintf("{%v %v}", runners.DomriaPrimary, runners.DomriaSecondary)
}
