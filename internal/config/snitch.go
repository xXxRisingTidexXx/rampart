package config

type Snitch struct {
	DSN    string `yaml:"-"`
	Spec   string `yaml:"spec"`
	Server Server `yaml:"server"`
}
