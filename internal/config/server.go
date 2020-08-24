package config

import (
	"fmt"
	"time"
)

type Server struct {
	Address        string        `yaml:"address"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
}

func (server *Server) String() string {
	return fmt.Sprintf(
		"{%s %s %s %d}",
		server.Address,
		server.ReadTimeout,
		server.WriteTimeout,
		server.MaxHeaderBytes,
	)
}
