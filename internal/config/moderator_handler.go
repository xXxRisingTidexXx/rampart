package config

type ModeratorHandler struct {
	Admin             string `yaml:"-"`
	StartCommand      string `yaml:"start-command"`
	ImageMarkupButton string `yaml:"image-markup-button"`
	Separator         string `yaml:"separator"`
}
