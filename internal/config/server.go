package config

import (
	"fmt"
	"time"
)

type Server struct {
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
}

func (server *Server) String() string {
	return fmt.Sprintf("{%s %s %d}", server.ReadTimeout, server.WriteTimeout, server.MaxHeaderBytes)
}
