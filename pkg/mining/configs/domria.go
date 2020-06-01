package configs

import (
	"fmt"
)

type Domria struct {
	Fetcher   *Fetcher   `yaml:"fetcher"`
	Validator *Validator `yaml:"validator"`
}

func (domria *Domria) String() string {
	return fmt.Sprintf("{%v %v}", domria.Fetcher, domria.Validator)
}
