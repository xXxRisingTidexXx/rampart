package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func NewConfig() (Config, error) {
	var config Config
	config.Messis.DSN = os.Getenv("RAMPART_DATABASE_DSN")
	if config.Messis.DSN == "" {
		return config, fmt.Errorf("config: failed to find the db dsn")
	}
	config.Telegram.Token = os.Getenv("RAMPART_TELEGRAM_TOKEN")
	if config.Telegram.Token == "" {
		return config, fmt.Errorf("config: failed to find the telegram token")
	}
	bytes, err := ioutil.ReadFile(misc.ResolvePath("config/dev.yaml"))
	if err != nil {
		return config, fmt.Errorf("config: failed to read the config file, %v", err)
	}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config, fmt.Errorf("config: failed to unmarshal the config file, %v", err)
	}
	config.Warhol.InputPath = misc.ResolvePath(config.Warhol.InputPath)
	config.Telegram.DSN = config.Messis.DSN
	return config, nil
}

type Config struct {
	Messis   Messis   `yaml:"messis"`
	Warhol   Warhol   `yaml:"warhol"`
	Telegram Telegram `yaml:"telegram"`
}
