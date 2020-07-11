package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Fetcher struct {
	Timeout   Timing                  `yaml:"timeout"`
	Portion   int                     `yaml:"portion"`
	Flags     map[misc.Housing]string `yaml:"flags"`
	Headers   map[string]string       `yaml:"headers"`
	SearchURL string                  `yaml:"searchURL"`
	SRID      int                     `yaml:"srid"`
}

func (fetcher *Fetcher) String() string {
	return fmt.Sprintf(
		"{%s %d %v %v %s %d}",
		fetcher.Timeout,
		fetcher.Portion,
		fetcher.Flags,
		fetcher.Headers,
		fetcher.SearchURL,
		fetcher.SRID,
	)
}
