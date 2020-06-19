package migrations

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"rampart/internal/homedir"
)

func newMigrations() (*migrations, error) {
	bytes, err := ioutil.ReadFile(homedir.Resolve("config/migrations.yaml"))
	if err != nil {
		return nil, fmt.Errorf("migrations: failed to read the config file, %v", err)
	}
	var migrations migrations
	if err = yaml.Unmarshal(bytes, &migrations); err != nil {
		return nil, fmt.Errorf("migrations: failed to unmarshal the config file, %v", err)
	}
	return &migrations, nil
}

type migrations struct {
	QueryParams map[string]string `yaml:"queryParams"`
}

func (migrations *migrations) String() string {
	return fmt.Sprintf("{%v}", migrations.QueryParams)
}
