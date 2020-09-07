package config

import (
	"fmt"
)

type Mining struct {
	DomriaPrimaryMiner   *DomriaMiner `yaml:"domriaPrimaryMiner"`
	DomriaSecondaryMiner *DomriaMiner `yaml:"domriaSecondaryMiner"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%v %v}", mining.DomriaPrimaryMiner, mining.DomriaSecondaryMiner)
}
