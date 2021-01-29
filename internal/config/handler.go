package config

type Handler struct {
	StartCommand               string `yaml:"start-command"`
	StartButton                string `yaml:"start-button"`
	HelpCommand                string `yaml:"help-command"`
	HelpButton                 string `yaml:"help-button"`
	CancelButton               string `yaml:"cancel-button"`
	AddButton                  string `yaml:"add-button"`
	ListButton                 string `yaml:"list-button"`
	AnyPriceButton             string `yaml:"any-price-button"`
	AnyRoomNumberButton        string `yaml:"any-room-number-button"`
	OneRoomNumberButton        string `yaml:"one-room-number-button"`
	TwoRoomNumberButton        string `yaml:"two-room-number-button"`
	ThreeRoomNumberButton      string `yaml:"three-room-number-button"`
	ManyRoomNumberButton       string `yaml:"many-room-number-button"`
	AnyFloorButton             string `yaml:"any-floor-button"`
	LowFloorButton             string `yaml:"low-floor-button"`
	HighFloorButton            string `yaml:"high-floor-button"`
	TemplatePath               string `yaml:"template-path"` // TODO: use it in helper's randomization.
	MinFlatCount               int    `yaml:"min-flat-count"`
	MaxPriceLength             int    `yaml:"max-price-length"`
	MaxRoomNumberLength        int    `yaml:"max-room-number-length"`
	MaxRoomNumber              int64  `yaml:"max-room-number"`
	AnyPricePlaceholder        string `yaml:"any-price-placeholder"`
	AnyRoomNumberPlaceholder   string `yaml:"any-room-number-placeholder"`
	OneRoomNumberPlaceholder   string `yaml:"one-room-number-placeholder"`
	TwoRoomNumberPlaceholder   string `yaml:"two-room-number-placeholder"`
	ThreeRoomNumberPlaceholder string `yaml:"three-room-number-placeholder"`
	ManyRoomNumberPlaceholder  string `yaml:"many-room-number-placeholder"`
	AnyFloorPlaceholder        string `yaml:"any-floor-placeholder"`
	LowFloorPlaceholder        string `yaml:"low-floor-placeholder"`
	HighFloorPlaceholder       string `yaml:"high-floor-placeholder"`
	DeleteButton               string `yaml:"delete-button"`
	DeletePrefix               string `yaml:"delete-prefix"`
}
