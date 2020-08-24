package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Gauger struct {
	InterpreterURL string       `yaml:"interpreterURL"`
	Headers        misc.Headers `yaml:"headers"`
	SearchRadius   float64      `yaml:"searchRadius"`
	MinArea        float64      `yaml:"minArea"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf(
		"{%s %v %f %f}",
		gauger.InterpreterURL,
		gauger.Headers,
		gauger.SearchRadius,
		gauger.MinArea,
	)
}
