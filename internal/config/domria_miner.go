package config

import (
	"time"
)

type DomriaMiner struct {
	Name         string        `yaml:"name"`
	Spec         string        `yaml:"spec"`
	Timeout      time.Duration `yaml:"timeout"`
	RetryLimit   int           `yaml:"retry-limit"`
	SearchPrefix string        `yaml:"search-prefix"`
	UserAgent    string        `yaml:"user-agent"`
}
