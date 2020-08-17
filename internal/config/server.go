package config

import (
	"fmt"
	"time"
)

type Server struct {
	Port           string        `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
}

func (server *Server) String() string {
	return fmt.Sprintf(
		"{%s %s %s %d}",
		server.Port,
		server.ReadTimeout,
		server.WriteTimeout,
		server.MaxHeaderBytes,
	)
}
