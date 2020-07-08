package config

import (
	"fmt"
)

type Verdicts struct {
	Valid   string `yaml:"valid"`
	Invalid string `yaml:"invalid"`
}

func (verdicts *Verdicts) String() string {
	return fmt.Sprintf("{%s %s}", verdicts.Valid, verdicts.Invalid)
}
