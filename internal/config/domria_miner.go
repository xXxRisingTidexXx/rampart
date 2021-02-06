package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type DomriaMiner struct {
	Alias     string       `yaml:"alias"`
	Housing   misc.Housing `yaml:"housing"`
	Spec      string       `yaml:"spec"`
	Server    Server       `yaml:"server"`
	Fetcher   Fetcher      `yaml:"fetcher"`
	Sanitizer Sanitizer    `yaml:"sanitizer"`
	Geocoder  Geocoder     `yaml:"geocoder"`
	Gauger    Gauger       `yaml:"gauger"`
	Validator Validator    `yaml:"validator"`
	Storer    Storer       `yaml:"storer"`
}

func (m DomriaMiner) Name() string {
	return m.Alias
}

func (m DomriaMiner) Schedule() string {
	return m.Spec
}

func (m DomriaMiner) Metrics() Server {
	return m.Server
}
