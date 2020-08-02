package secrets

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/homedir"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func NewSecrets() (*Secrets, error) {
	bytes, err := ioutil.ReadFile(homedir.Resolve("secrets/dev.yaml"))
	if err != nil {
		return nil, fmt.Errorf("secrets: failed to read the secrets file, %v", err)
	}
	secrets := Secrets{}
	if err = yaml.Unmarshal(bytes, &secrets); err != nil {
		return nil, fmt.Errorf("secrets: failed to unmarshal the secrets file, %v", err)
	}
	return &secrets, nil
}

type Secrets struct {
	DSN string `yaml:"dsn"`
}

func (secrets *Secrets) String() string {
	return fmt.Sprintf("{%s}", secrets.DSN)
}
