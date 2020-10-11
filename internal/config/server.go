package config

import (
	"time"
)

type Server struct {
	Address        string        `yaml:"address"`
	ReadTimeout    time.Duration `yaml:"read-timeout"`
	WriteTimeout   time.Duration `yaml:"write-timeout"`
	MaxHeaderBytes int           `yaml:"max-header-bytes"`
}
