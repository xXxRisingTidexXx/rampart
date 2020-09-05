package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Gauger struct {
	InterpreterURL string       `yaml:"interpreterURL"`
	Headers        misc.Headers `yaml:"headers"`
	SearchRadius   float64      `yaml:"searchRadius"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf("{%s %v %f}", gauger.InterpreterURL, gauger.Headers, gauger.SearchRadius)
}
