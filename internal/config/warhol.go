package config

import (
	"time"
)

type Warhol struct {
	Timeout      time.Duration `yaml:"timeout"`
	ThreadNumber int           `yaml:"thread-number"`
	RetryLimit   int           `yaml:"retry-limit"`
	OutputFormat string        `yaml:"output-format"`
	InputPath    string        `yaml:"input-path"`
}
