package config

import (
	"fmt"
	"rampart/internal/misc"
)

type Server struct {
	Port int `yaml:"port"`
	ReadTimeout misc.Timing `yaml:"readTimeout"`
	WriteTimeout misc.Timing `yaml:"writeTimeout"`
	MaxHeaderBytes int `yaml:"maxHeaderBytes"`
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
