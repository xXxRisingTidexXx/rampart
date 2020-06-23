package config

import (
	"fmt"
)

type Miners struct {
	DomriaPrimary   *Domria `yaml:"domriaPrimary"`
	DomriaSecondary *Domria `yaml:"domriaSecondary"`
}

func (miners *Miners) String() string {
	return fmt.Sprintf("{%v %v}", miners.DomriaPrimary, miners.DomriaSecondary)
}
