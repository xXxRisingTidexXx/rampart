package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Fetcher struct {
	Timeout   Timing                  `yaml:"timeout"`
	Portion   int                     `yaml:"portion"`
	Flags     map[misc.Housing]string `yaml:"flags"`
	Headers   map[string]string       `yaml:"headers"`
	SearchURL string                  `yaml:"searchURL"`
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
