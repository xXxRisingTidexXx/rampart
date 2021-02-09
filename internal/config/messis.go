package config

type Messis struct {
	DSN                string             `yaml:"-"`
	DomriaMiner        DomriaMiner        `yaml:"domria-miner"`
	GeocodingAmplifier GeocodingAmplifier `yaml:"geocoding-amplifier"`
	Server             Server             `yaml:"server"`
}
