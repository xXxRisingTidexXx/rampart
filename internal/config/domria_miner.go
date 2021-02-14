package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type DomriaMiner struct {
	Name                    string            `yaml:"name"`
	Spec                    string            `yaml:"spec"`
	Timeout                 time.Duration     `yaml:"timeout"`
	Page                    int               `yaml:"page"`
	RetryLimit              int               `yaml:"retry-limit"`
	SearchPrefix            string            `yaml:"search-prefix"`
	UserAgent               string            `yaml:"user-agent"`
	URLPrefix               string            `yaml:"url-prefix"`
	ImageURLFormat          string            `yaml:"image-url-format"`
	MaxTotalArea            float64           `yaml:"max-total-area"`
	MaxRoomNumber           int               `yaml:"max-room-number"`
	MaxTotalFloor           int               `yaml:"max-total-floor"`
	Swaps                   misc.Set          `yaml:"swaps"`
	Cities                  map[string]string `yaml:"cities"`
	StreetReplacements      []string          `yaml:"street-replacements"`
	HouseNumberReplacements []string          `yaml:"house-number-replacements"`
	MaxHouseNumberLength    int               `yaml:"max-house-number-length"`
}
