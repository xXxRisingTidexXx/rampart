package config

import (
	"fmt"
)

type Verdict struct {
	Valid   string `yaml:"valid"`
	Invalid string `yaml:"invalid"`
}

func (verdict *Verdict) String() string {
	return fmt.Sprintf("{%s %s}", verdict.Valid, verdict.Invalid)
}
