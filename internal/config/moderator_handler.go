package config

type ModeratorHandler struct {
	Admin              string `yaml:"-"`
	StartCommand       string `yaml:"start-command"`
	MarkupButton       string `yaml:"markup-button"`
	LuxuryButton       string `yaml:"luxury-button"`
	ComfortButton      string `yaml:"comfort-button"`
	JunkButton         string `yaml:"junk-button"`
	ConstructionButton string `yaml:"construction-button"`
	ExcessButton       string `yaml:"excess-button"`
	EnoughButton       string `yaml:"enough-button"`
}
