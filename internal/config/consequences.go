package config

import (
	"fmt"
)

type Consequences struct {
	Created   string `yaml:"created"`
	Updated   string `yaml:"updated"`
	Unaltered string `yaml:"unaltered"`
	Failed    string `yaml:"failed"`
}

func (consequences *Consequences) String() string {
	return fmt.Sprintf(
		"{%s %s %s %s}",
		consequences.Created,
		consequences.Updated,
		consequences.Unaltered,
		consequences.Failed,
	)
}
