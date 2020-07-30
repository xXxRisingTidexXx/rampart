package config

import (
	"fmt"
)

type Storer struct {
	SRID int `yaml:"srid"`
}

func (storer *Storer) String() string {
	return fmt.Sprintf("{%d}", storer.SRID)
}
