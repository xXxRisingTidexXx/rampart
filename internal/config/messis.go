package config

type Messis struct {
	DSN                  string      `yaml:"-"`
	DomriaPrimaryMiner   DomriaMiner `yaml:"domria-primary-miner"`
	DomriaSecondaryMiner DomriaMiner `yaml:"domria-secondary-miner"`
}
