package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Fetcher struct {
	Timeout   time.Duration      `yaml:"timeout"`
	Portion   int                `yaml:"portion"`
	Flags     map[Housing]string `yaml:"flags"`
	Headers   misc.Headers       `yaml:"headers"`
	SearchURL string             `yaml:"search-url"`
}
