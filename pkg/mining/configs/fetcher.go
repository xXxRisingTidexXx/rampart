package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"rampart/pkg/mining"
	"time"
)

type Fetcher struct {
	Timeout         time.Duration
	Portion         int
	Flags           map[mining.Housing]string
	SearchURL       string
	OriginURL       string
	ImageURL        string
	USDLabel        string
	MinTotalArea    float64
	MaxTotalArea    float64
	MinLivingArea   float64
	MinKitchenArea  float64
	MinRoomNumber   int
	MaxRoomNumber   int
	MinSpecificArea float64
	MaxSpecificArea float64
	MinFloor        int
	MinTotalFloor   int
	MaxTotalFloor   int
	MinLongitude    float64
	MaxLongitude    float64
	MinLatitude     float64
	MaxLatitude     float64
	StateEnding     string
	StateSuffix     string
	DistrictLabel   string
	DistrictEnding  string
	DistrictSuffix  string
}

func (fetcher *Fetcher) UnmarshalYAML(node *yaml.Node) error {
	type Alias struct {
		Timeout         string                    `yaml:"timeout"`
		Portion         int                       `yaml:"portion"`
		Flags           map[mining.Housing]string `yaml:"flags"`
		SearchURL       string                    `yaml:"searchURL"`
		OriginURL       string                    `yaml:"originURL"`
		ImageURL        string                    `yaml:"imageURL"`
		USDLabel        string                    `yaml:"usdLabel"`
		MinTotalArea    float64                   `yaml:"minTotalArea"`
		MaxTotalArea    float64                   `yaml:"maxTotalArea"`
		MinLivingArea   float64                   `yaml:"minLivingArea"`
		MinKitchenArea  float64                   `yaml:"minKitchenArea"`
		MinRoomNumber   int                       `yaml:"minRoomNumber"`
		MaxRoomNumber   int                       `yaml:"maxRoomNumber"`
		MinSpecificArea float64                   `yaml:"minSpecificArea"`
		MaxSpecificArea float64                   `yaml:"maxSpecificArea"`
		MinFloor        int                       `yaml:"minFloor"`
		MinTotalFloor   int                       `yaml:"minTotalFloor"`
		MaxTotalFloor   int                       `yaml:"maxTotalFloor"`
		MinLongitude    float64                   `yaml:"minLongitude"`
		MaxLongitude    float64                   `yaml:"maxLongitude"`
		MinLatitude     float64                   `yaml:"minLatitude"`
		MaxLatitude     float64                   `yaml:"maxLatitude"`
		StateEnding     string                    `yaml:"stateEnding"`
		StateSuffix     string                    `yaml:"stateSuffix"`
		DistrictLabel   string                    `yaml:"districtLabel"`
		DistrictEnding  string                    `yaml:"districtEnding"`
		DistrictSuffix  string                    `yaml:"districtSuffix"`
	}
	var alias Alias
	if err := node.Decode(&alias); err != nil {
		return err
	}
	timeout, err := time.ParseDuration(alias.Timeout)
	if err != nil {
		return err
	}
	fetcher.Timeout = timeout
	fetcher.Portion = alias.Portion
	fetcher.Flags = alias.Flags
	fetcher.SearchURL = alias.SearchURL
	fetcher.OriginURL = alias.OriginURL
	fetcher.ImageURL = alias.ImageURL
	fetcher.USDLabel = alias.USDLabel
	fetcher.MinTotalArea = alias.MinTotalArea
	fetcher.MaxTotalArea = alias.MaxTotalArea
	fetcher.MinLivingArea = alias.MinLivingArea
	fetcher.MinKitchenArea = alias.MinKitchenArea
	fetcher.MinRoomNumber = alias.MinRoomNumber
	fetcher.MaxRoomNumber = alias.MaxRoomNumber
	fetcher.MinSpecificArea = alias.MinSpecificArea
	fetcher.MaxSpecificArea = alias.MaxSpecificArea
	fetcher.MinFloor = alias.MinFloor
	fetcher.MinTotalFloor = alias.MinTotalFloor
	fetcher.MaxTotalFloor = alias.MaxTotalFloor
	fetcher.MinLongitude = alias.MinLongitude
	fetcher.MaxLongitude = alias.MaxLongitude
	fetcher.MinLatitude = alias.MinLatitude
	fetcher.MaxLatitude = alias.MaxLatitude
	fetcher.StateEnding = alias.StateEnding
	fetcher.StateSuffix = alias.StateSuffix
	fetcher.DistrictLabel = alias.DistrictLabel
	fetcher.DistrictEnding = alias.DistrictEnding
	fetcher.DistrictSuffix = alias.DistrictSuffix
	return nil
}

func (fetcher *Fetcher) String() string {
	return fmt.Sprintf(
		"{%s %d %v %s %s %s %s %.1f %.1f %.1f %.1f %d %d %.1f "+
			"%.1f %d %d %d %.1f %.1f %.1f %.1f %s %s %s %s %s}",
		fetcher.Timeout,
		fetcher.Portion,
		fetcher.Flags,
		fetcher.SearchURL,
		fetcher.OriginURL,
		fetcher.ImageURL,
		fetcher.USDLabel,
		fetcher.MinTotalArea,
		fetcher.MaxTotalArea,
		fetcher.MinLivingArea,
		fetcher.MinKitchenArea,
		fetcher.MinRoomNumber,
		fetcher.MaxRoomNumber,
		fetcher.MinSpecificArea,
		fetcher.MaxSpecificArea,
		fetcher.MinFloor,
		fetcher.MinTotalFloor,
		fetcher.MaxTotalFloor,
		fetcher.MinLongitude,
		fetcher.MaxLongitude,
		fetcher.MinLatitude,
		fetcher.MaxLatitude,
		fetcher.StateEnding,
		fetcher.StateSuffix,
		fetcher.DistrictLabel,
		fetcher.DistrictEnding,
		fetcher.DistrictSuffix,
	)
}
