package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Fetcher struct {
	Timeout      time.Duration           `yaml:"timeout"`
	RetryLimit   int                     `yaml:"retry-limit"`
	Portion      int                     `yaml:"portion"`
	Flags        map[misc.Housing]string `yaml:"flags"`
	SearchFormat string                  `yaml:"search-format"`
}
