package config

import (
	"fmt"
)

type Storer struct {
	UpdateTiming Timing `yaml:"updateTiming"`
}

func (storer *Storer) String() string {
	return fmt.Sprintf("{%s}", storer.UpdateTiming)
}
