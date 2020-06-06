package configs

import (
	"fmt"
)

type Domria struct {
	Fetcher   *Fetcher   `yaml:"fetcher"`
	Sanitizer *Sanitizer `yaml:"sanitizer"`
	Validator *Validator `yaml:"validator"`
}

func (domria *Domria) String() string {
	return fmt.Sprintf("{%v %v %v}", domria.Fetcher, domria.Sanitizer, domria.Validator)
}
