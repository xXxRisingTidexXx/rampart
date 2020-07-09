package misc

import (
	"fmt"
)

type Verdicts struct {
	Valid   string `yaml:"valid"`
	Invalid string `yaml:"invalid"`
}

func (verdicts *Verdicts) Targets() []string {
	return []string{verdicts.Valid, verdicts.Invalid}
}

func (verdicts *Verdicts) String() string {
	return fmt.Sprintf("{%s %s}", verdicts.Valid, verdicts.Invalid)
}
