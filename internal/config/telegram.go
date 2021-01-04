package config

type Telegram struct {
	Token  string `yaml:"-"`
	DSN    string `yaml:"-"`
	Server Server `yaml:"server"`
}
