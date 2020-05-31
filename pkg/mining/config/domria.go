package config

import (
	"fmt"
)

type Domria struct {
	Fetcher *Fetcher `yaml:"fetcher"`
}

func (domria *Domria) String() string {
	return fmt.Sprintf("{%v}", domria.Fetcher)
}
