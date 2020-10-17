package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func NewConfig() (Config, error) {
	var config Config
	bytes, err := ioutil.ReadFile(misc.ResolvePath("config/dev.yaml"))
	if err != nil {
		return config, fmt.Errorf("config: failed to read the config file, %v", err)
	}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config, fmt.Errorf("config: failed to unmarshal the config file, %v", err)
	}
	return config, nil
}

type Config struct {
	Mining Mining `yaml:"mining"`
}
