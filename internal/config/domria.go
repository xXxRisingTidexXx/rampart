package config

import (
	"fmt"
)

type Domria struct {
	Fetcher   *Fetcher   `yaml:"fetcher"`
	Sanitizer *Sanitizer `yaml:"sanitizer"`
	Geocoder  *Geocoder  `yaml:"geocoder"`
	Validator *Validator `yaml:"validator"`
	Sifter    *Sifter    `yaml:"sifter"`
}

func (domria *Domria) String() string {
	return fmt.Sprintf(
		"{%v %v %v %v %v}",
		domria.Fetcher,
		domria.Sanitizer,
		domria.Geocoder,
		domria.Validator,
		domria.Sifter,
	)
}
