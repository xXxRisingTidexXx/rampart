package config

type Telegram struct {
	Token   string `yaml:"-"`
	DSN     string `yaml:"-"`
	Timeout int    `yaml:"timeout"`
	Server  Server `yaml:"server"`
}
