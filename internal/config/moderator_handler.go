package config

type ModeratorHandler struct {
	Admin              string `yaml:"-"`
	StartCommand       string `yaml:"start-command"`
	StartButton        string `yaml:"start-button"`
	ImageMarkupButton  string `yaml:"image-markup-button"`
	LuxuryButton       string `yaml:"luxury-button"`
	ComfortButton      string `yaml:"comfort-button"`
	JunkButton         string `yaml:"junk-button"`
	ConstructionButton string `yaml:"construction-button"`
	ExcessButton       string `yaml:"excess-button"`
}
