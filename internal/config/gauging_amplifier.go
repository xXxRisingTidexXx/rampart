package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"time"
)

type GaugingAmplifier struct {
	Timeout           time.Duration `yaml:"timeout"`
	Host              string        `yaml:"host"`
	InterpreterFormat string        `yaml:"interpreter-format"`
	SubwayCities      misc.Set      `yaml:"subway-cities"`
	SSFSearchRadius   float64       `yaml:"ssf-search-radius"`
	SSFMinDistance    float64       `yaml:"ssf-min-distance"`
	SSFModifier       float64       `yaml:"ssf-modifier"`
	IZFSearchRadius   float64       `yaml:"izf-search-radius"`
	IZFMinArea        float64       `yaml:"izf-min-area"`
	IZFMinDistance    float64       `yaml:"izf-min-distance"`
	IZFModifier       float64       `yaml:"izf-modifier"`
	GZFSearchRadius   float64       `yaml:"gzf-search-radius"`
	GZFMinArea        float64       `yaml:"gzf-min-area"`
	GZFMinDistance    float64       `yaml:"gzf-min-distance"`
	GZFModifier       float64       `yaml:"gzf-modifier"`
}
