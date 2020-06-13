package configs

import (
	"fmt"
	"rampart/internal/mining"
	"rampart/internal/mining/util"
)

type Fetcher struct {
	Timeout   util.Timeout              `yaml:"timeout"`
	Portion   int                       `yaml:"portion"`
	Flags     map[mining.Housing]string `yaml:"flags"`
	Headers   map[string]string         `yaml:"headers"`
	SearchURL string                    `yaml:"searchURL"`
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
