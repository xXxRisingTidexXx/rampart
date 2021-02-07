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
	MaxTotalArea           float64           `yaml:"max-total-area"`
	MaxRoomNumber          float64           `yaml:"max-room-number"`
	MaxTotalFloor          float64           `yaml:"max-total-floor"`
	Swaps                  misc.Set          `yaml:"swaps"`
	CityOrthography        map[string]string `yaml:"city-orthography"`
	StreetOrthography      []string          `yaml:"street-orthography"`
	HouseNumberOrthography []string          `yaml:"house-number-orthography"`
	MaxHouseNumberLength   int               `yaml:"max-house-number-length"`
}
