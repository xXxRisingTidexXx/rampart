package config

import (
	"fmt"
)

type Statuses struct {
	Success string `yaml:"success"`
	Failure string `yaml:"failure"`
}

func (statuses *Statuses) String() string {
	return fmt.Sprintf("{%s %s}", statuses.Success, statuses.Failure)
}
