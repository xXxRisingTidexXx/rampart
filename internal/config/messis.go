package config

type Messis struct {
	DSN                string             `yaml:"-"`
	DomriaMiner        DomriaMiner        `yaml:"domria-miner"`
	BufferSize         int                `yaml:"buffer-size"`
	GeocodingAmplifier GeocodingAmplifier `yaml:"geocoding-amplifier"`
	GaugingAmplifiers  []GaugingAmplifier `yaml:"gauging-amplifiers"`
	StoringAmplifier   StoringAmplifier   `yaml:"storing-amplifier"`
	Server             Server             `yaml:"server"`
}
