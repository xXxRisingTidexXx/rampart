package misc

import (
	"fmt"
)

type Statuses struct {
	Success string `yaml:"success"`
	Failure string `yaml:"failure"`
}

func (statuses *Statuses) Targets() []string {
	return []string{statuses.Success, statuses.Failure}
}

func (statuses *Statuses) String() string {
	return fmt.Sprintf("{%s %s}", statuses.Success, statuses.Failure)
}
