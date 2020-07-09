package misc

import (
	"fmt"
)

type Verdicts struct {
	Approved string `yaml:"approved"`
	Denied   string `yaml:"denied"`
}

func (verdicts *Verdicts) Targets() []string {
	return []string{verdicts.Approved, verdicts.Denied}
}

func (verdicts *Verdicts) String() string {
	return fmt.Sprintf("{%s %s}", verdicts.Approved, verdicts.Denied)
}
