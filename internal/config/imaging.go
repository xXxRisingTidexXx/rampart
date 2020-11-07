package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Imaging struct {
	ThreadNumber int          `yaml:"thread-number"`
	Headers      misc.Headers `yaml:"headers"`
	OutputFormat string       `yaml:"output-format"`
	InputPath    string       `yaml:"input-string"`
}
