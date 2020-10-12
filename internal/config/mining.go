package config

type Mining struct {
	DomriaPrimaryMiner   DomriaMiner `yaml:"domria-primary-miner"`
	DomriaSecondaryMiner DomriaMiner `yaml:"domria-secondary-miner"`
}
