package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"rampart/internal/mining"
	"time"
)

type Fetcher struct {
	Timeout   time.Duration
	Portion   int
	Flags     map[mining.Housing]string
	Headers   map[string]string
	SearchURL string
}

func (fetcher *Fetcher) UnmarshalYAML(node *yaml.Node) error {
	type Alias struct {
		Timeout   string                    `yaml:"timeout"`
		Portion   int                       `yaml:"portion"`
		Flags     map[mining.Housing]string `yaml:"flags"`
		Headers   map[string]string         `yaml:"headers"`
		SearchURL string                    `yaml:"searchURL"`
	}
	var alias Alias
	if err := node.Decode(&alias); err != nil {
		return err
	}
	timeout, err := time.ParseDuration(alias.Timeout)
	if err != nil {
		return err
	}
	fetcher.Timeout = timeout
	fetcher.Portion = alias.Portion
	fetcher.Flags = alias.Flags
	fetcher.Headers = alias.Headers
	fetcher.SearchURL = alias.SearchURL
	return nil
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
