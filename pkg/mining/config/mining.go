package config

import (
	"fmt"
)

type Mining struct {
	UserAgent   string       `yaml:"userAgent"`
	Prospectors *Prospectors `yaml:"prospectors"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%s %v}", mining.UserAgent, mining.Prospectors)
}
