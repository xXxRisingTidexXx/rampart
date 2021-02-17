package config

type Assistant struct {
	DSN        string              `yaml:"-"`
	Token      string              `yaml:"-"`
	Publisher  Publisher           `yaml:"publisher"`
	Spec       string              `yaml:"spec"`
	Server     Server              `yaml:"server"`
	Dispatcher AssistantDispatcher `yaml:"dispatcher"`
}