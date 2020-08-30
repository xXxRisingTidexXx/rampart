package domria

import (
	"encoding/json"
	"fmt"
)

type item struct {
	BeautifulURL        string
	MainPhoto           string
	Photos              map[string]*photo
	Panoramas           []*panorama
	UpdatedAt           moment
	PriceArr            *prices
	TotalSquareMeters   float64
	LivingSquareMeters  float64
	KitchenSquareMeters float64
	RoomsCount          int
	Floor               int
	FloorsCount         int
	UserNewbuildNameUK  string
	Longitude           coordinate
	Latitude            coordinate
	StateNameUK         string
	CityNameUK          string
	DistrictNameUK      string
	StreetNameUK        string
	StreetName          string
	BuildingNumberStr   string
	Source              string
}

func (item *item) UnmarshalJSON(bytes []byte) error {
	sourceless := sourcelessItem{}
	if err := json.Unmarshal(bytes, &sourceless); err != nil {
		return err
	}
	item.BeautifulURL = sourceless.BeautifulURL
	item.MainPhoto = sourceless.MainPhoto
	item.Photos = sourceless.Photos
	item.Panoramas = sourceless.Panoramas
	item.UpdatedAt = sourceless.UpdatedAt
	item.PriceArr = sourceless.PriceArr
	item.TotalSquareMeters = sourceless.TotalSquareMeters
	item.LivingSquareMeters = sourceless.LivingSquareMeters
	item.KitchenSquareMeters = sourceless.KitchenSquareMeters
	item.RoomsCount = sourceless.RoomsCount
	item.Floor = sourceless.Floor
	item.FloorsCount = sourceless.FloorsCount
	item.UserNewbuildNameUK = sourceless.UserNewbuildNameUK
	item.Longitude = sourceless.Longitude
	item.Latitude = sourceless.Latitude
	item.StateNameUK = sourceless.StateNameUK
	item.CityNameUK = sourceless.CityNameUK
	item.DistrictNameUK = sourceless.DistrictNameUK
	item.StreetNameUK = sourceless.StreetNameUK
	item.StreetName = sourceless.StreetName
	item.BuildingNumberStr = sourceless.BuildingNumberStr
	item.Source = string(bytes)
	return nil
}

func (item *item) String() string {
	return fmt.Sprintf(
		"{%s %s %v %v %s %v %.1f %.1f %.1f %d %d %d %s %.6f %.6f %s %s %s %s %s %s %s}",
		item.BeautifulURL,
		item.MainPhoto,
		item.Photos,
		item.Panoramas,
		item.UpdatedAt,
		item.PriceArr,
		item.TotalSquareMeters,
		item.LivingSquareMeters,
		item.KitchenSquareMeters,
		item.RoomsCount,
		item.Floor,
		item.FloorsCount,
		item.UserNewbuildNameUK,
		item.Longitude,
		item.Latitude,
		item.StateNameUK,
		item.CityNameUK,
		item.DistrictNameUK,
		item.StreetNameUK,
		item.StreetName,
		item.BuildingNumberStr,
		item.Source,
	)
}
