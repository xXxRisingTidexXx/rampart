package config

type Snitch struct {
	DSN       string    `yaml:"-"`
	Publisher Publisher `yaml:"publisher"`
	Spec      string    `yaml:"spec"`
	Server    Server    `yaml:"server"`
}
