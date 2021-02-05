package config

type Telegram struct {
	DSN        string     `yaml:"-"`
	Token      string     `yaml:"-"`
	Dispatcher Dispatcher `yaml:"dispatcher"`
	Publisher  Publisher  `yaml:"publisher"`
	Spec       string     `yaml:"spec"`
	Server     Server     `yaml:"server"`
}
