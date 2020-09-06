package config

import (
	"fmt"
)

type DomriaMiner struct {
	Alias     string     `yaml:"alias"`
	Housing   Housing    `yaml:"housing"`
	Spec      string     `yaml:"spec"`
	Server    *Server    `yaml:"server"`
	Fetcher   *Fetcher   `yaml:"fetcher"`
	Sanitizer *Sanitizer `yaml:"sanitizer"`
	Geocoder  *Geocoder  `yaml:"geocoder"`
	Gauger    *Gauger    `yaml:"gauger"`
	Validator *Validator `yaml:"validator"`
	Storer    *Storer    `yaml:"storer"`
}

func (miner *DomriaMiner) Name() string {
	return miner.Alias
}

func (miner *DomriaMiner) Schedule() string {
	return miner.Spec
}

func (miner *DomriaMiner) Metrics() *Server {
	return miner.Server
}

func (miner *DomriaMiner) String() string {
	return fmt.Sprintf(
		"{%s %s %s %v %v %v %v %v %v %v}",
		miner.Alias,
		miner.Housing,
		miner.Spec,
		miner.Server,
		miner.Fetcher,
		miner.Sanitizer,
		miner.Geocoder,
		miner.Gauger,
		miner.Validator,
		miner.Storer,
	)
}
