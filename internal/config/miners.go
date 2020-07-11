package config

import (
	"fmt"
)

type Miners struct {
	DomriaPrimary   *DomriaMiner `yaml:"domriaPrimary"`
	DomriaSecondary *DomriaMiner `yaml:"domriaSecondary"`
}

func (miners *Miners) String() string {
	return fmt.Sprintf("{%v %v}", miners.DomriaPrimary, miners.DomriaSecondary)
}
