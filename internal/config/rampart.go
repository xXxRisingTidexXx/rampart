package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"rampart/internal/homedir"
)

func NewRampart() (*Rampart, error) {
	bytes, err := ioutil.ReadFile(homedir.Resolve("config/dev.yaml"))
	if err != nil {
		return nil, fmt.Errorf("config: rampart failed to read the config file, %v", err)
	}
	var rampart Rampart
	if err = yaml.Unmarshal(bytes, &rampart); err != nil {
		return nil, fmt.Errorf("config: rampart failed to unmarshal the config file, %v", err)
	}
	return &rampart, nil
}

type Rampart struct {
	Mining     *Mining     `yaml:"mining"`
	Migrations *Migrations `yaml:"migrations"`
}

func (rampart *Rampart) String() string {
	return fmt.Sprintf("{%v %v}", rampart.Mining, rampart.Migrations)
}
