package config

import (
	"fmt"
)

type State struct {
	Ending string `yaml:"ending"`
	Suffix string `yaml:"suffix"`
}

func (state *State) String() string {
	return fmt.Sprintf("{%s %s}", state.Ending, state.Suffix)
}
