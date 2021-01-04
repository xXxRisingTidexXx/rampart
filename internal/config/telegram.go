package config

type Telegram struct {
	DSN        string     `yaml:"-"`
	Dispatcher Dispatcher `yaml:"dispatcher"`
	Server     Server     `yaml:"server"`
}
