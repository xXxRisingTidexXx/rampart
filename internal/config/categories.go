package config

import (
	"fmt"
)

type Categories struct {
	Located      string `yaml:"located"`
	Unlocated    string `yaml:"unlocated"`
	Failed       string `yaml:"failed"`
	Inconclusive string `yaml:"inconclusive"`
	Successful   string `yaml:"successful"`
}

func (categories *Categories) String() string {
	return fmt.Sprintf(
		"{%s %s %s %s %s}",
		categories.Located,
		categories.Unlocated,
		categories.Failed,
		categories.Inconclusive,
		categories.Successful,
	)
}
