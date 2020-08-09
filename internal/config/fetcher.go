package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type Fetcher struct {
	Timeout   time.Duration      `yaml:"timeout"`
	Portion   int                `yaml:"portion"`
	Flags     map[Housing]string `yaml:"flags"`
	Headers   misc.Headers       `yaml:"headers"`
	SearchURL string             `yaml:"searchURL"`
}

func (fetcher *Fetcher) String() string {
	return fmt.Sprintf(
		"{%s %d %v %v %s}",
		fetcher.Timeout,
		fetcher.Portion,
		fetcher.Flags,
		fetcher.Headers,
		fetcher.SearchURL,
	)
}
