package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"rampart/pkg/mining"
	"time"
)

type Fetcher struct {
	Timeout         time.Duration
	Portion         int
	Flags           map[mining.Housing]string
	SearchURL       string
	OriginURLPrefix string
	ImageURLPrefix  string
	StateEnding     string
	StateSuffix     string
	DistrictLabel   string
	DistrictEnding  string
	DistrictSuffix  string
}

func (fetcher *Fetcher) UnmarshalYAML(node *yaml.Node) error {
	type Alias struct {
		Timeout         string                    `yaml:"timeout"`
		Portion         int                       `yaml:"portion"`
		Flags           map[mining.Housing]string `yaml:"flags"`
		SearchURL       string                    `yaml:"searchURL"`
		OriginURLPrefix string                    `yaml:"originURLPrefix"`
		ImageURLPrefix  string                    `yaml:"imageURLPrefix"`
		StateEnding     string                    `yaml:"stateEnding"`
		StateSuffix     string                    `yaml:"stateSuffix"`
		DistrictLabel   string                    `yaml:"districtLabel"`
		DistrictEnding  string                    `yaml:"districtEnding"`
		DistrictSuffix  string                    `yaml:"districtSuffix"`
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
	fetcher.SearchURL = alias.SearchURL
	fetcher.OriginURLPrefix = alias.OriginURLPrefix
	fetcher.ImageURLPrefix = alias.ImageURLPrefix
	fetcher.StateEnding = alias.StateEnding
	fetcher.StateSuffix = alias.StateSuffix
	fetcher.DistrictLabel = alias.DistrictLabel
	fetcher.DistrictEnding = alias.DistrictEnding
	fetcher.DistrictSuffix = alias.DistrictSuffix
	return nil
}

func (fetcher *Fetcher) String() string {
	return fmt.Sprintf(
		"{%s %d %v %s %s %s %s %s %s %s %s}",
		fetcher.Timeout,
		fetcher.Portion,
		fetcher.Flags,
		fetcher.SearchURL,
		fetcher.OriginURLPrefix,
		fetcher.ImageURLPrefix,
		fetcher.StateEnding,
		fetcher.StateSuffix,
		fetcher.DistrictLabel,
		fetcher.DistrictEnding,
		fetcher.DistrictSuffix,
	)
}
