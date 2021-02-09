package config

type Messis struct {
	DSN                string             `yaml:"-"`
	DomriaMiner        DomriaMiner        `yaml:"domria-miner"`
	GeocodingAmplifier GeocodingAmplifier `yaml:"geocoding-amplifier"`
	GaugingAmplifiers  []GaugingAmplifier `yaml:"gauging-amplifiers"`
	Server             Server             `yaml:"server"`
}
