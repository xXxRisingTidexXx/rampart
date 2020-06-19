package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"rampart/internal/homedir"
)

func NewMining() (*Mining, error) {
	bytes, err := ioutil.ReadFile(homedir.Resolve("config/mining.yaml"))
	if err != nil {
		return nil, fmt.Errorf("config: mining failed to read the config file, %v", err)
	}
	var mining Mining
	if err = yaml.Unmarshal(bytes, &mining); err != nil {
		return nil, fmt.Errorf("config: mining failed to unmarshal the config file, %v", err)
	}
	return &mining, nil
}

type Mining struct {
	SRID        int          `yaml:"srid"`
	Prospectors *Prospectors `yaml:"prospectors"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%d %v}", mining.SRID, mining.Prospectors)
}
