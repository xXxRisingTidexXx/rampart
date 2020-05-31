package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func NewMining() (*Mining, error) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("config: failed to find the caller path")
	}
	file, err := os.Open(
		filepath.Join(
			filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filePath)))),
			"config",
			"mining.yaml",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("config: failed to open the config file, %v", err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read the config file, %v", err)
	}
	if err = file.Close(); err != nil {
		return nil, fmt.Errorf("config: failed to close the config file, %v", err)
	}
	var mining Mining
	if err = yaml.Unmarshal(bytes, &mining); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal the config file, %v", err)
	}
	return &mining, nil
}

type Mining struct {
	UserAgent   string       `yaml:"userAgent"`
	Prospectors *Prospectors `yaml:"prospectors"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%s %v}", mining.UserAgent, mining.Prospectors)
}
