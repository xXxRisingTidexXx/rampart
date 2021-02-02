package config

type Snitch struct {
	DSN    string `yaml:"-"`
	Server Server `yaml:"server"`
}
