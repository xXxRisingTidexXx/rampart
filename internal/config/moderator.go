package config

type Moderator struct {
	DSN   string `yaml:"-"`
	Token string `yaml:"-"`
	Admin string `yaml:"-"`
}
