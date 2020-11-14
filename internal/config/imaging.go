package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Imaging struct {
	Timeout      time.Duration `yaml:"timeout"`
	ThreadNumber int           `yaml:"thread-number"`
	RetryLimit   int           `yaml:"retry-limit"`
	Headers      misc.Headers  `yaml:"headers"`
	OutputFormat string        `yaml:"output-format"`
	InputPath    string        `yaml:"input-path"`
}
