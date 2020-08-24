package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/homedir"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func NewConfig() (*Config, error) {
	bytes, err := ioutil.ReadFile(homedir.Resolve("config/dev.yaml"))
	if err != nil {
		return nil, fmt.Errorf("config: failed to read the config file, %v", err)
	}
	config := Config{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal the config file, %v", err)
	}
	return &config, nil
}

type Config struct {
	Mining     *Mining     `yaml:"mining"`
	Migrations *Migrations `yaml:"migrations"`
	Gauging    *Gauging    `yaml:"gauging"`
}

func (config *Config) String() string {
	return fmt.Sprintf("{%v %v %v}", config.Mining, config.Migrations, config.Gauging)
}
