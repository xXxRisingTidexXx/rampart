package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Gauger struct {
	Timeout        time.Duration `yaml:"timeout"`
	Headers        misc.Headers  `yaml:"headers"`
	InterpreterURL string        `yaml:"interpreterURL"`
	SubwayCities   misc.Set      `yaml:"subwayCities"`
}

func (gauger *Gauger) String() string {
	return fmt.Sprintf(
		"{%s %v %s %v}",
		gauger.Timeout,
		gauger.Headers,
		gauger.InterpreterURL,
		gauger.SubwayCities,
	)
}
