package config

import (
	"fmt"
	"rampart/internal/misc"
)

type DomriaMiner struct {
	Alias     string       `yaml:"alias"`
	Housing   misc.Housing `yaml:"housing"`
	Spec      string       `yaml:"spec"`
	Port      int          `yaml:"port"`
	Fetcher   *Fetcher     `yaml:"fetcher"`
	Sanitizer *Sanitizer   `yaml:"sanitizer"`
	Geocoder  *Geocoder    `yaml:"geocoder"`
	Validator *Validator   `yaml:"validator"`
	Storer    *Storer      `yaml:"storer"`
}

func (domria *DomriaMiner) String() string {
	return fmt.Sprintf(
		"{%s %s %s %d %v %v %v %v %v}",
		domria.Alias,
		domria.Housing,
		domria.Spec,
		domria.Port,
		domria.Fetcher,
		domria.Sanitizer,
		domria.Geocoder,
		domria.Validator,
		domria.Storer,
	)
}
