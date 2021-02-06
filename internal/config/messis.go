package config

type Messis struct {
	DSN         string      `yaml:"-"`
	DomriaMiner DomriaMiner `yaml:"domria-miner"`
	Server      Server      `yaml:"server"`
}
