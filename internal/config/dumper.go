package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Dumper struct {
	Timeout    time.Duration `yaml:"timeout"`
	GaugingURL string        `yaml:"gaugingURL"`
	Headers    misc.Headers  `yaml:"headers"`
}

func (dumper *Dumper) String() string {
	return fmt.Sprintf("{%s %s %v}", dumper.Timeout, dumper.GaugingURL, dumper.Headers)
}
