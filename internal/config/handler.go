package config

type Handler struct {
	StartCommand          string `yaml:"start-command"`
	StartButton           string `yaml:"start-button"`
	HelpCommand           string `yaml:"help-command"`
	HelpButton            string `yaml:"help-button"`
	CancelButton          string `yaml:"cancel-button"`
	AddButton             string `yaml:"add-button"`
	AnyPriceButton        string `yaml:"any-price-button"`
	AnyRoomNumberButton   string `yaml:"any-room-number-button"`
	OneRoomNumberButton   string `yaml:"one-room-number-button"`
	TwoRoomNumberButton   string `yaml:"two-room-number-button"`
	ThreeRoomNumberButton string `yaml:"three-room-number-button"`
	ManyRoomNumberButton  string `yaml:"many-room-number-button"`
	AnyFloorButton        string `yaml:"any-floor-button"`
	LowFloorButton        string `yaml:"low-floor-button"`
	HighFloorButton       string `yaml:"high-floor-button"`
	TemplatePath          string `yaml:"template-path"`  // TODO: use it in helper's randomization.
	MinFlatCount          int    `yaml:"min-flat-count"`
	MaxRoomNumber         int64  `yaml:"max-room-number"`
}
