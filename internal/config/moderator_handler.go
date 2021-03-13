package config

type ModeratorHandler struct {
	StartCommand       string   `yaml:"start-command"`
	StartButton        string   `yaml:"start-button"`
	HelpCommand        string   `yaml:"help-command"`
	HelpButton         string   `yaml:"help-button"`
	MarkupButton       string   `yaml:"markup-button"`
	LuxuryButton       string   `yaml:"luxury-button"`
	ComfortButton      string   `yaml:"comfort-button"`
	JunkButton         string   `yaml:"junk-button"`
	ConstructionButton string   `yaml:"construction-button"`
	ExcessButton       string   `yaml:"excess-button"`
	EnoughButton       string   `yaml:"enough-button"`
}
