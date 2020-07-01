package config

import (
	"fmt"
)

type Server struct {
	Port           int    `yaml:"port"`
	ReadTimeout    Timing `yaml:"readTimeout"`
	WriteTimeout   Timing `yaml:"writeTimeout"`
	MaxHeaderBytes int    `yaml:"maxHeaderBytes"`
}

func (server *Server) String() string {
	return fmt.Sprintf(
		"{%d %s %s %d}",
		server.Port,
		server.ReadTimeout,
		server.WriteTimeout,
		server.MaxHeaderBytes,
	)
}
