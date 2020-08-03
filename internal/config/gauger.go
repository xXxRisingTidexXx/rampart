package config

import (
	"fmt"
)

type Gauger struct {
	Timeout        Timing            `yaml:"timeout"`
	Headers        map[string]string `yaml:"headers"`
	InterpreterURL string            `yaml:"interpreterURL"`
	SearchRadius   float64           `yaml:"searchRadius"`
	NoDistance     float64           `yaml:"noDistance"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf(
		"{%s %v %s %.1f %f}",
		gauger.Timeout,
		gauger.Headers,
		gauger.InterpreterURL,
		gauger.SearchRadius,
		gauger.NoDistance,
	)
}
