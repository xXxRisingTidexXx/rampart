package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Storer struct {
	UpdateTiming   misc.Timing `yaml:"updateTiming"`
}

func (storer *Storer) String() string {
	return fmt.Sprintf("{%s}", storer.UpdateTiming)
}
