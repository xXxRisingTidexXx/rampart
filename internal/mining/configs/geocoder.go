package configs

import (
	"fmt"
	"rampart/internal/mining/util"
)

type Geocoder struct {
	Timeout   util.Timeout      `yaml:"timeout"`
	Headers   map[string]string `yaml:"headers"`
	Stateless *util.Set         `yaml:"stateless"`
	SearchURL string            `yaml:"searchURL"`
	MinLookup float64           `yaml:"minLookup"`
}

func (geocoder *Geocoder) String() string {
	return fmt.Sprintf(
		"{%s %v %v %s %.2f}",
		geocoder.Timeout,
		geocoder.Headers,
		geocoder.Stateless,
		geocoder.SearchURL,
		geocoder.MinLookup,
	)
}
