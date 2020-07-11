package config

import (
	"fmt"
)

type Server struct {
	ReadTimeout    Timing `yaml:"readTimeout"`
	WriteTimeout   Timing `yaml:"writeTimeout"`
	MaxHeaderBytes int    `yaml:"maxHeaderBytes"`
}

func (server *Server) String() string {
	return fmt.Sprintf("{%s %s %d}", server.ReadTimeout, server.WriteTimeout, server.MaxHeaderBytes)
}
