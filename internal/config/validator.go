package config

type Validator struct {
	MinMediaCount   int     `yaml:"min-media-count"`
	MinPrice        float64 `yaml:"min-price"`
	MinTotalArea    float64 `yaml:"min-total-area"`
	MaxTotalArea    float64 `yaml:"max-total-area"`
	MinLivingArea   float64 `yaml:"min-living-area"`
	MinKitchenArea  float64 `yaml:"min-kitchen-area"`
	MinRoomNumber   int     `yaml:"min-room-number"`
	MaxRoomNumber   int     `yaml:"max-room-number"`
	MinSpecificArea float64 `yaml:"min-specific-area"`
	MaxSpecificArea float64 `yaml:"max-specific-area"`
	MinFloor        int     `yaml:"min-floor"`
	MinTotalFloor   int     `yaml:"min-total-floor"`
	MaxTotalFloor   int     `yaml:"max-total-floor"`
	MinLongitude    float64 `yaml:"min-longitude"`
	MaxLongitude    float64 `yaml:"max-longitude"`
	MinLatitude     float64 `yaml:"min-latitude"`
	MaxLatitude     float64 `yaml:"max-latitude"`
}
