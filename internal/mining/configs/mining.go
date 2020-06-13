package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func NewMining() (mining *Mining, err error) {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("configs: mining failed to find the caller path")
	}
	file, err := os.Open(
		filepath.Join(
			filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filePath)))),
			"configs",
			"mining.yaml",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("configs: mining failed to open the config file, %v", err)
	}
	defer func() {
		if err == nil {
			if err = file.Close(); err != nil {
				mining = nil
				err = fmt.Errorf("configs: mining failed to close the config file, %v", err)
			}
		}
	}()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("configs: mining failed to read the config file, %v", err)
	}
	var config Mining
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("configs: mining failed to unmarshal the config file, %v", err)
	}
	return &config, nil
}

type Mining struct {
	SRID        int          `yaml:"srid"`
	Prospectors *Prospectors `yaml:"prospectors"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf("{%d %v}", mining.SRID, mining.Prospectors)
}
