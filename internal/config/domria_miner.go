package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type DomriaMiner struct {
	Name                   string            `yaml:"name"`
	Spec                   string            `yaml:"spec"`
	Timeout                time.Duration     `yaml:"timeout"`
	RetryLimit             int               `yaml:"retry-limit"`
	SearchPrefix           string            `yaml:"search-prefix"`
	UserAgent              string            `yaml:"user-agent"`
	URLPrefix              string            `yaml:"url-prefix"`
	ImageURLFormat         string            `yaml:"image-url-format"`
	Swaps                  misc.Set          `yaml:"swaps"`
	CityOrthography        map[string]string `yaml:"city-orthography"`
	StreetOrthography      []string          `yaml:"street-orthography"`
	HouseNumberOrthography []string          `yaml:"house-number-orthography"`
	HouseNumberMaxLength   int               `yaml:"house-number-max-length"`
}
