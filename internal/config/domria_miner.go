package config

import (
	"fmt"
)

type DomriaMiner struct {
	Alias     string     `yaml:"alias"`
	Housing   Housing    `yaml:"housing"`
	Spec      string     `yaml:"spec"`
	Port      int        `yaml:"port"`
	Fetcher   *Fetcher   `yaml:"fetcher"`
	Sanitizer *Sanitizer `yaml:"sanitizer"`
	Geocoder  *Geocoder  `yaml:"geocoder"`
	Gauger    *Gauger    `yaml:"gauger"`
	Validator *Validator `yaml:"validator"`
	Storer    *Storer    `yaml:"storer"`
}

func (domria *DomriaMiner) String() string {
	return fmt.Sprintf(
		"{%s %s %s %d %v %v %v %v %v %v}",
		domria.Alias,
		domria.Housing,
		domria.Spec,
		domria.Port,
		domria.Fetcher,
		domria.Sanitizer,
		domria.Geocoder,
		domria.Gauger,
		domria.Validator,
		domria.Storer,
	)
}
