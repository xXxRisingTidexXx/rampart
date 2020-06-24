package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Domria struct {
	Alias     string
	Housing   misc.Housing `yaml:"housing"`
	Spec      string       `yaml:"spec"`
	Fetcher   *Fetcher     `yaml:"fetcher"`
	Sanitizer *Sanitizer   `yaml:"sanitizer"`
	Geocoder  *Geocoder    `yaml:"geocoder"`
	Validator *Validator   `yaml:"validator"`
	Storer    *Storer      `yaml:"storer"`
}

func (domria *Domria) String() string {
	return fmt.Sprintf(
		"{%s %s %s %v %v %v %v %v}",
		domria.Alias,
		domria.Housing,
		domria.Spec,
		domria.Fetcher,
		domria.Sanitizer,
		domria.Geocoder,
		domria.Validator,
		domria.Storer,
	)
}
