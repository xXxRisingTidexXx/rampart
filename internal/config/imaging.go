package config

import (
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

type Imaging struct {
	InputPath    string       `yaml:"input-string"`
	ThreadNumber int          `yaml:"thread-number"`
	Headers      misc.Headers `yaml:"headers"`
	OutputFormat string       `yaml:"output-format"`
}
