package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Domria struct {
	Fetcher   *Fetcher     `yaml:"fetcher"`
	Sanitizer *Sanitizer   `yaml:"sanitizer"`
	Geocoder  *Geocoder    `yaml:"geocoder"`
	Validator *Validator   `yaml:"validator"`
	Storer    *Storer      `yaml:"storer"`
	Housing   misc.Housing `yaml:"housing"`
	Spec      string       `yaml:"spec"`
}

func (domria *Domria) String() string {
	return fmt.Sprintf(
		"{%v %v %v %v %v %s %s}",
		domria.Fetcher,
		domria.Sanitizer,
		domria.Geocoder,
		domria.Validator,
		domria.Storer,
		domria.Housing,
		domria.Spec,
	)
}
