package config

type Moderator struct {
	DSN        string              `yaml:"-"`
	Token      string              `yaml:"-"`
	Dispatcher ModeratorDispatcher `yaml:"dispatcher"`
}
