package config

import (
	"fmt"
)

type District struct {
	Label  string `yaml:"label"`
	Ending string `yaml:"ending"`
	Suffix string `yaml:"suffix"`
}

func (district *District) String() string {
	return fmt.Sprintf("{%s %s %s}", district.Label, district.Ending, district.Suffix)
}
