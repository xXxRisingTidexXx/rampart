package config

type ModeratorDispatcher struct {
	Timeout int              `yaml:"timeout"`
	Handler ModeratorHandler `yaml:"handler"`
}
